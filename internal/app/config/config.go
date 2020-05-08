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

// LibraryConfig specifies the library settings
type LibraryConfig struct {
	BorrowExpireDays uint `json:"borrow_expire_days"`
	DdlExtendDays    uint `json:"ddl_extend_days"`
	MaxExtendTimes   uint `json:"max_extend_times"`
	MaxOverdueBooks  uint `json:"max_overdue_books"`
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

// LoadLibraryConfig reads the library settings from the file
func LoadLibraryConfig(file string) (LibraryConfig, error) {
	var cfg LibraryConfig
	cfgFile, err := os.Open(file)
	defer cfgFile.Close()
	if err != nil {
		fmt.Println("[error] LoadLibraryConfig: unable to open the config file.")
		return cfg, err
	}
	dec := json.NewDecoder(cfgFile)
	if err = dec.Decode(&cfg); err != nil {
		fmt.Println("[error] LoadLibraryConfig: invalid configuration.")
		return cfg, err
	}
	return cfg, nil
}
