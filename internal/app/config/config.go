package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// DbConfig specifies the database connection settings
type DbConfig struct {
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// LoadDbConfig reads the database connection settings from the file
func LoadDbConfig(file string) (DbConfig, error) {
	var cfg DbConfig
	cfgFile, err := os.Open(file)
	defer cfgFile.Close()
	if err != nil {
		fmt.Println("[error] LoadDbConfig: unable to open the config file.")
		return cfg, err
	}
	dec := json.NewDecoder(cfgFile)
	if err = dec.Decode(&cfg); err != nil {
		fmt.Println("[error] LoadDbConfig: invalid configuration.")
		return cfg, err
	}
	return cfg, nil
}
