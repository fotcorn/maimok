package maimok

import (
	"net/http"

	"github.com/go-chi/render"
)

// ListVMsHandler handler
func (state *globalState) ListVMsHandler(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, vm := range ListVMs(state) {
		list = append(list, vm)
	}
	render.RenderList(w, r, list)
}

// CreateVMHandler handles POST / requests
func (state *globalState) CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	createVM := CreateVMStruct{
		DiskSpaceGB: 100,
		RAMMB:       1024,
		Name:        "TestVM",
		Image:       "xenial-server-cloudimg-amd64-disk1.img",
	}
	CreateVM(state, createVM)
}

// Render a VM
func (rd *VM) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
