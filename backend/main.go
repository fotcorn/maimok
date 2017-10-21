package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/vms", func(r chi.Router) {
		r.Get("/", ListVMs)
	})

	http.ListenAndServe(":3000", r)
}

// ListVMs handler
func ListVMs(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, vm := range vms {
		list = append(list, vm)
	}
	render.RenderList(w, r, list)
}

// Render a VM
func (rd *VM) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// VM data model
type VM struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Memory int    `json:"memory"`
}

var vms = []*VM{
	{ID: "1234", Name: "1_nextcloud", Memory: 1024},
	{ID: "4321", Name: "2_test", Memory: 512},
}
