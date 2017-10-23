package maimok

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// GetRouter returns the router configuration
func GetRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/vms", func(r chi.Router) {
		r.Get("/", ListVMsHandler)
		r.Post("/", CreateVMHandler)
	})
	return r
}
