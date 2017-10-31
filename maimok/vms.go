package maimok

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	libvirt "github.com/libvirt/libvirt-go"
	"github.com/satori/go.uuid"
)

// VM data model
type VM struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Memory uint64  `json:"memory"`
	State  VMState `json:"state"`
}

// VMState is the state of the vm
type VMState int

const (
	// Running vm
	Running VMState = iota + 1
	// Stopped vm
	Stopped
)

// CreateVMStruct struct for the CreateVM call
type CreateVMStruct struct {
	ID          string
	Name        string
	RAMMB       uint
	DiskSpaceGB uint
	Image       string
	SSHKey      string
}

// CreateVM createds a virtual machine
func CreateVM(state *globalState, createVM CreateVMStruct) error {
	createVM.ID = uuid.NewV4().String()
	createVM.SSHKey = state.config.SSHKey

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
	var out bytes.Buffer
	cmd.Stdout = &out
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
	len, err := stream.Send(buf.Bytes())
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

	// create domain
	/*buf = new(bytes.Buffer)
	if err = state.tpl.ExecuteTemplate(buf, "templates/domain.xml", createVM); err != nil {
		return err
	}
	if _, err = state.conn.DomainCreateXML(buf.String(), libvirt.DOMAIN_NONE); err != nil {
		return err
	}*/
	return nil
}

// ListVMs returns a list of all virtual machines
func ListVMs(state *globalState) []*VM {
	domains, err := state.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		return nil
	}

	vms := []*VM{}
	for _, domain := range domains {
		name, _ := domain.GetName()
		memory, _ := domain.GetMaxMemory()
		id, _ := domain.GetID()
		state, _, _ := domain.GetState()

		var vmState VMState
		if state == libvirt.DOMAIN_RUNNING {
			vmState = Running
		} else {
			vmState = Stopped
		}

		vms = append(vms, &VM{ID: id, Name: name, Memory: memory, State: vmState})
	}
	return vms
}
