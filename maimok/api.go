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

type createVMResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (resp *createVMResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// CreateVMHandler handles POST / requests
func (state *globalState) CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	createVM := CreateVMStruct{
		DiskSpaceGB: 100,
		RAMMB:       1024,
		Name:        "TestVM",
		Image:       "xenial-server-cloudimg-amd64-disk1.img",
	}
	err := CreateVM(state, createVM)
	if err != nil {
		render.Render(w, r, &createVMResponse{Status: "error", Error: err.Error()})
	} else {
		render.Render(w, r, &createVMResponse{Status: "ok"})
	}
}

// Render a VM
func (rd *VM) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
