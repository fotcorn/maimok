package maimok

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config is the
type Config struct {
	LibvirtURL string   `toml:"libvirt_url"`
	Image      string   `toml:"image"`
	IPRange    []string `toml:"ip_range"`
}

// LoadConfig load Config from file
func LoadConfig() (*Config, error) {
	config := Config{}
	_, err := toml.DecodeFile("config.toml", &config)
	fmt.Println(config)
	if err != nil {
		return &Config{}, err
	}
	return &config, nil
}
