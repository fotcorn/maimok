package maimok

import (
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

// CreateVM createds a virtual machine
func CreateVM(name string, ram uint, ipAddress string) string {
	return "Creating VM"
	// create storage
	// create image
	//conn, err := libvirt.NewConnect("qemu+ssh://sysadmin@10.0.0.200/system")

}

// ListVMs test
func ListVMs() []*VM {
	conn, err := libvirt.NewConnect("qemu+ssh://sysadmin@10.0.0.200/system")
	if err != nil {
		return nil
	}
	defer conn.Close()

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
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
