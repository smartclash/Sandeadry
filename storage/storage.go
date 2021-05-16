package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StorageInit() (err error) {
	_, err = gorm.Open(sqlite.Open(""), &gorm.Config{})

	return
}
