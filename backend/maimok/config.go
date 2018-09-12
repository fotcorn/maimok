package maimok

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config is the structure of the config file on disk
type Config struct {
	LibvirtURL string `toml:"libvirt_url"`
	Image      string `toml:"image"`
	SSHKey     string `toml:"ssh_key"`
	Gateway    string `toml:"gateway"`
	Netmask    string `toml:"netmask"`
}

// LoadConfig load Config from file
func LoadConfig() (*Config, error) {
	config := Config{}
	md, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		return &Config{}, err
	}
	if len(md.Undecoded()) > 0 {
		return &Config{}, fmt.Errorf("Invalid config parameter: %q", md.Undecoded())
	}
	if !md.IsDefined("libvirt_url") {
		config.LibvirtURL = "qemu:///system"
	}
	if !md.IsDefined("image") {
		return &Config{}, fmt.Errorf("\"image\" parameter is required")
	}
	if !md.IsDefined("gateway") {
		return &Config{}, fmt.Errorf("\"gateway\" parameter is required")
	}
	if !md.IsDefined("netmask") {
		return &Config{}, fmt.Errorf("\"netmask\" parameter is required")
	}
	return &config, nil
}
