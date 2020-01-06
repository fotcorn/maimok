package maimok

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	libvirt "github.com/libvirt/libvirt-go"
	uuid "github.com/satori/go.uuid"
)

// VM data model
type VM struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Memory  uint64 `json:"memory"`
	Running bool   `json:"running"`
}

// CreateVMStruct struct for the CreateVM call
type CreateVMStruct struct {
	// Generated
	ID         string
	MACAddress string
	// Provided by API
	Name        string
	Hostname    string
	RAMMB       uint
	DiskSpaceGB uint
	Image       string
	IPAddress   string
	// Set in config file
	SSHKey  string
	Gateway string
	Netmask string
}

func generateMACAddress() (string, error) {
	buf := make(net.HardwareAddr, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	buf[0] = (buf[0] | 2) & 0xfe // Set local bit, ensure unicast address
	return buf.String(), nil
}

// CreateVM createds a virtual machine
func CreateVM(state *globalState, createVM CreateVMStruct) error {
	createVM.ID = uuid.NewV4().String()
	createVM.SSHKey = state.config.SSHKey
	createVM.Gateway = state.config.Gateway
	createVM.Netmask = state.config.Netmask

	macAddress, err := generateMACAddress()
	if err != nil {
		return fmt.Errorf("Cannot generate MAC Address, %s", err)
	}
	createVM.MACAddress = macAddress

	// create config iso
	dir, err := ioutil.TempDir("", "maimok")
	if err != nil {
		return fmt.Errorf("Cannot create temp directory, %s", err)
	}
	defer os.RemoveAll(dir)

	metaDataFile, err := os.OpenFile(filepath.Join(dir, "meta-data"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer metaDataFile.Close()
	if err := state.tpl.ExecuteTemplate(metaDataFile, "meta-data.yml", createVM); err != nil {
		return err
	}

	userDataFile, err := os.OpenFile(filepath.Join(dir, "user-data"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer userDataFile.Close()
	if err := state.tpl.ExecuteTemplate(userDataFile, "user-data.yml", createVM); err != nil {
		return err
	}

	networkConfigFile, err := os.OpenFile(filepath.Join(dir, "network-config"), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer networkConfigFile.Close()
	if err := state.tpl.ExecuteTemplate(networkConfigFile, "network-config.yml", createVM); err != nil {
		return err
	}

	cmd := exec.Command("genisoimage", "-volid", "cidata", "-joliet", "-rock",
		metaDataFile.Name(), userDataFile.Name(), networkConfigFile.Name())
	var isoFile bytes.Buffer
	cmd.Stdout = &isoFile
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Generation of iso image failed: %s", err)
	}

	storagePool, err := state.conn.LookupStoragePoolByName("default")
	if err != nil {
		return err
	}

	// create iso volume
	buf := new(bytes.Buffer)
	if err := state.tpl.ExecuteTemplate(buf, "iso-volume.xml", createVM); err != nil {
		return err
	}
	volume, err := storagePool.StorageVolCreateXML(buf.String(), 0)
	if err != nil {
		return err
	}

	stream, err := state.conn.NewStream(0)
	if err != nil {
		return err
	}
	volume.Upload(stream, 0, 0, 0)
	len, err := stream.Send(isoFile.Bytes())
	if err != nil {
		return fmt.Errorf("Could not send config iso image to storage pool: %s", err)
	}
	if len < buf.Len() {
		return fmt.Errorf("Could not send all iso image data to storage pool")
	}
	if err := stream.Finish(); err != nil {
		return fmt.Errorf("Could not send config iso image to storage pool: %s", err)
	}

	// create harddisk volume
	buf = new(bytes.Buffer)
	if err := state.tpl.ExecuteTemplate(buf, "volume.xml", createVM); err != nil {
		return err
	}
	if _, err = storagePool.StorageVolCreateXML(buf.String(), 0); err != nil {
		return err
	}

	// define domain
	buf = new(bytes.Buffer)
	if err = state.tpl.ExecuteTemplate(buf, "domain.xml", createVM); err != nil {
		return err
	}
	domain, err := state.conn.DomainDefineXML(buf.String())
	if err != nil {
		return fmt.Errorf("Failed to define domain: %s", err)
	}

	// start domain
	if err = domain.Create(); err != nil {
		return fmt.Errorf("Failed to start domain: %s", err)
	}

	return nil
}

// ListVMs returns a list of all virtual machines
func ListVMs(state *globalState) []*VM {
	domains, err := state.conn.ListAllDomains(0)
	if err != nil {
		return nil
	}

	vms := []*VM{}
	for _, domain := range domains {
		name, _ := domain.GetName()
		memory, _ := domain.GetMaxMemory()
		id, _ := domain.GetID()
		state, _, _ := domain.GetState()

		var running bool
		if state == libvirt.DOMAIN_RUNNING {
			running = true
		} else {
			running = false
		}

		vms = append(vms, &VM{ID: id, Name: name, Memory: memory, Running: running})
	}

	sort.Slice(vms, func(i, j int) bool {
		return vms[i].Name < vms[j].Name
	})

	return vms
}

// StartVM starts a virtual machine with the given name
func StartVM(state *globalState, name string) error {
	domain, err := state.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	domain.Create()
	return nil
}

// StopVM tries to gracefully stop a VM, which might fail
func StopVM(state *globalState, name string) error {
	domain, err := state.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	domain.Shutdown()
	return nil
}

// ForceStopVM forcefully stops a VM, like pressing the power button on a real PC
func ForceStopVM(state *globalState, name string) error {
	domain, err := state.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	domain.Destroy()
	return nil
}

// RestartVM tries to gracefully restart a VM, which might fail
func RestartVM(state *globalState, name string) error {
	domain, err := state.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	domain.Reboot(0)
	return nil
}

// ForceRestartVM forcefully restarts/resets a VM, like pressing the reset button on a real PC
func ForceRestartVM(state *globalState, name string) error {
	domain, err := state.conn.LookupDomainByName(name)
	if err != nil {
		return err
	}
	domain.Reset(0)
	return nil
}
