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

// Render a VM
func (rd *VM) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
