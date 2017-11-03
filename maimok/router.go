package maimok

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// GetRouter returns the router configuration
func GetRouter(state *globalState) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/vms", func(r chi.Router) {
		r.Get("/", state.ListVMsHandler)
		r.Post("/", state.CreateVMHandler)
	})
	return r
}
