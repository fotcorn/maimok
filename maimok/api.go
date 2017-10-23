package maimok

import (
	"net/http"

	"github.com/go-chi/render"
)

// ListVMsHandler handler
func ListVMsHandler(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, vm := range ListVMs() {
		list = append(list, vm)
	}
	render.RenderList(w, r, list)
}

// CreateVMHandler handles POST / requests
func CreateVMHandler(w http.ResponseWriter, r *http.Request) {

}

// Render a VM
func (rd *VM) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
