package maimok

import (
	"net/http"

	"github.com/go-chi/render"
)

func validationFailed(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, 400)
	render.Render(w, r, &createVMResponse{Status: "error", Error: message})
	return
}

// ListVMsHandler handler
func (state *globalState) ListVMsHandler(w http.ResponseWriter, r *http.Request) {
	list := []render.Renderer{}
	for _, vm := range ListVMs(state) {
		list = append(list, vm)
	}
	render.RenderList(w, r, list)
}

type createVMRequest struct {
	DiskSpaceGB uint   `json:"disk_space_gb"`
	RAMMB       uint   `json:"ram_mb"`
	Name        string `json:"name"`
	Hostname    string `json:"hostname"`
	Image       string `json:"image"`
	IPAddress   string `json:"ip_address"`
}

func (a *createVMRequest) Bind(r *http.Request) error {
	return nil
}

type createVMResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (resp *createVMResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// CreateVMHandler handles POST / requests
func (state *globalState) CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	request := &createVMRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, &createVMResponse{Status: "error", Error: err.Error()})
		return
	}

	if request.DiskSpaceGB == 0 {
		validationFailed(w, r, "disk_space_gb field is required")
		return
	}
	if request.RAMMB == 0 {
		validationFailed(w, r, "ram_mb field is required")
		return
	}
	if request.Image == "" {
		validationFailed(w, r, "image field is required")
		return
	}
	if request.IPAddress == "" {
		validationFailed(w, r, "ip_address field is required")
		return
	}
	if request.Name == "" {
		validationFailed(w, r, "name field is required")
		return
	}
	if request.Hostname == "" {
		validationFailed(w, r, "hostname field is required")
		return
	}

	createVM := CreateVMStruct{
		DiskSpaceGB: request.DiskSpaceGB,
		RAMMB:       request.RAMMB,
		Name:        request.Name,
		Hostname:    request.Hostname,
		Image:       request.Image,
		IPAddress:   request.IPAddress,
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
