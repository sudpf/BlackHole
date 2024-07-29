package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
	DB *gorm.DB
}

func (m *MySQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	m.DB = db
	return db, nil
}

func (m *MySQLDatabase) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MySQLDatabase) CreateTable(model interface{}) error {
	return m.DB.AutoMigrate(model)
}

func (m *MySQLDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := m.DB.Where(conditions).Find(model)
	return query, query.Error
}
