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
	IPAddress   string
	Image       string
}

// CreateVM createds a virtual machine
func CreateVM(state *globalState, createVM CreateVMStruct) {
	// create volume
	buf := new(bytes.Buffer)
	err := state.tpl.ExecuteTemplate(buf, "templates/volume.xml", createVM)
	if err != nil {
		panic(err)
	}
	storagePool, err := state.conn.LookupStoragePoolByName("default")
	if err != nil {
		panic(err)
	}
	_, err = storagePool.StorageVolCreateXML(buf.String(), 0)
	if err != nil {
		panic(err)
	}

	// create domain
	buf = new(bytes.Buffer)
	err = state.tpl.ExecuteTemplate(buf, "templates/domain.xml", createVM)
	if err != nil {
		panic(err)
	}
	_, err = state.conn.DomainCreateXML(buf.String(), libvirt.DOMAIN_NONE)
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
