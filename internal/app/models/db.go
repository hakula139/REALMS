package models

import (
	"fmt"

	"github.com/hakula139/REALMS/internal/app/config"
	"github.com/jinzhu/gorm"
)

// DbSetup opens a database connection and initializes
func DbSetup(cfg config.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/?parseTime=true",
		cfg.Username, cfg.Password,
		cfg.Protocol, cfg.Host, cfg.Port)
	db, err := gorm.Open(cfg.Type, dsn)
	if err != nil {
		fmt.Println("[error] DbSetup: connection failed: " + err.Error())
		return nil, err
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.Database)
	db.Exec("USE " + cfg.Database)
	db.AutoMigrate(&Book{}, &User{}, &Record{})
	return db, nil
}
