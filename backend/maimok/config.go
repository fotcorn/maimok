package maimok

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the structure of the config file on disk
type Config struct {
	LibvirtURL string
	Image      string
	SSHKey     string
	Gateway    string
	Netmask    string
}

// LoadConfig loads configuration from file or environment variables
func LoadConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("libvirt_url", "qemu:///system")
	viper.SetEnvPrefix("maimok")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return &Config{}, err
	}

	config := Config{
		LibvirtURL: viper.GetString("libvirt_url"),
		Image:      viper.GetString("image"),
		SSHKey:     viper.GetString("ssh_key"),
		Gateway:    viper.GetString("gateway"),
		Netmask:    viper.GetString("netmask"),
	}

	if config.Image == "" {
		return &Config{}, fmt.Errorf("\"image\" config parameter is required")
	}
	if config.SSHKey == "" {
		return &Config{}, fmt.Errorf("\"ssh_key\" config parameter is required")
	}
	if config.Gateway == "" {
		return &Config{}, fmt.Errorf("\"gateway\" config parameter is required")
	}
	if config.Netmask == "" {
		return &Config{}, fmt.Errorf("\"netmask\" config parameter is required")
	}

	return &config, nil
}
