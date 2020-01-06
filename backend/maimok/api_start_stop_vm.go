package maimok

import (
	"net/http"

	"github.com/go-chi/render"
)

type request struct {
	Name string `json:"name"`
}

func (a *request) Bind(r *http.Request) error {
	return nil
}

type response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (resp *response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getName(w http.ResponseWriter, r *http.Request) string {
	req := &request{}
	if err := render.Bind(r, req); err != nil {
		render.Render(w, r, &response{Status: "error", Error: err.Error()})
		return ""
	}
	return req.Name
}

func sendResponse(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		render.Render(w, r, &response{Status: "error", Error: err.Error()})
	} else {
		render.Render(w, r, &response{Status: "ok"})
	}
}

func (state *globalState) StartVMHandler(w http.ResponseWriter, r *http.Request) {
	name := getName(w, r)
	if name != "" {
		err := StartVM(state, name)
		sendResponse(err, w, r)
	}
}
func (state *globalState) StopVMHandler(w http.ResponseWriter, r *http.Request) {
	name := getName(w, r)
	if name != "" {
		err := StopVM(state, name)
		sendResponse(err, w, r)
	}
}
func (state *globalState) ForceStopVMHandler(w http.ResponseWriter, r *http.Request) {
	name := getName(w, r)
	if name != "" {
		err := ForceStopVM(state, name)
		sendResponse(err, w, r)
	}
}
func (state *globalState) RestartVMHandler(w http.ResponseWriter, r *http.Request) {
	name := getName(w, r)
	if name != "" {
		err := RestartVM(state, name)
		sendResponse(err, w, r)
	}

}
func (state *globalState) ForceRestartVMHandler(w http.ResponseWriter, r *http.Request) {
	name := getName(w, r)
	if name != "" {
		err := ForceRestartVM(state, name)
		sendResponse(err, w, r)
	}

}
