package database

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDatabaseConnection = errors.New("Error - Falla de coneccion con la base de datos")
)

func errWraper(err error) error {
	return fmt.Errorf("%v: %v", ErrDatabaseConnection, err)
}

var (
	db   *gorm.DB
	once sync.Once
)

func GetDb() (_ *gorm.DB, err error) {
	once.Do(func() {
		dbUri := os.Getenv("DB_URI")
		db, err = gorm.Open(mysql.Open(dbUri), &gorm.Config{})
	})
	if err != nil {
		return nil, errWraper(err)
	}
	return db, err
}
