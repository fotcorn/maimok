package maimok

import (
	"bytes"

	libvirt "github.com/libvirt/libvirt-go"
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
	Name        string
	RAMMB       uint
	DiskSpaceGB uint
	Image       string
}

// CreateVM createds a virtual machine
func CreateVM(state *globalState, createVM CreateVMStruct) error {
	// create volume
	buf := new(bytes.Buffer)
	if err := state.tpl.ExecuteTemplate(buf, "templates/volume.xml", createVM); err != nil {
		return err
	}
	storagePool, err := state.conn.LookupStoragePoolByName("default")
	if err != nil {
		return err
	}
	if _, err = storagePool.StorageVolCreateXML(buf.String(), 0); err != nil {
		return err
	}

	// create domain
	buf = new(bytes.Buffer)
	if err = state.tpl.ExecuteTemplate(buf, "templates/domain.xml", createVM); err != nil {
		return err
	}
	if _, err = state.conn.DomainCreateXML(buf.String(), libvirt.DOMAIN_NONE); err != nil {
		return err
	}
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
