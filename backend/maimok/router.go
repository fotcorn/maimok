package maimok

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// GetRouter returns the router configuration
func GetRouter(state *globalState) chi.Router {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/vms", func(r chi.Router) {
		r.Get("/", state.ListVMsHandler)
		r.Post("/", state.CreateVMHandler)
	})

	r.Route("/api/vm", func(r chi.Router) {
		r.Post("/start", state.StartVMHandler)
		r.Post("/stop", state.StopVMHandler)
		r.Post("/forcestop", state.ForceStopVMHandler)
		r.Post("/restart", state.RestartVMHandler)
		r.Post("/forcerestart", state.ForceRestartVMHandler)
	})

	// serve dist/static/ directory with js/css/image/font files
	serveStaticFiles(r)

	// serve dist/index.html at root
	serveIndexHTML(r)

	return r
}

func serveStaticFiles(r chi.Router) {
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "dist/static")
	root := http.Dir(staticDir)

	path := "/static/"

	fs := http.StripPrefix(path, http.FileServer(root))

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func serveIndexHTML(r chi.Router) {
	workDir, _ := os.Getwd()
	indexFile := filepath.Join(workDir, "dist/index.html")
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexFile)
	}))
}
