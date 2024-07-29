package db

import "gorm.io/gorm"

type Database interface {
	Connect(connectionString string) (*gorm.DB, error)
	Close() error
	CreateTable(model interface{}) error
	Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error)
}
