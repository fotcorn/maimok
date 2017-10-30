package maimok

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	libvirt "github.com/libvirt/libvirt-go"
)

type globalState struct {
	config *Config
	conn   *libvirt.Connect
	tpl    *template.Template
}

// Run is the main entry point
func Run() {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	// load config file
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Cannot load configuration file: %s\n", err)
		code = 1
		return
	}

	// load templates
	templates, err := template.ParseFiles("templates/*")
	if err != nil {
		fmt.Printf("Cannot load template files: %s\n", err)
		code = 1
		return
	}

	// connect to hypervisor
	connection, err := libvirt.NewConnect(config.LibvirtURL)
	defer connection.Close()
	if err != nil {
		fmt.Printf("Cannot connect to libvirtd with URL %s, %s\n", config.LibvirtURL, err)
		code = 1
		return
	}
	state := globalState{conn: connection, config: config, tpl: templates}

	fmt.Println("Starting server on port 7000")
	err = http.ListenAndServe(":7000", GetRouter(&state))
	if err != nil {
		fmt.Println(err)
		code = 1
		return
	}
}
