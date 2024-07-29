package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDatabase struct {
	DB *gorm.DB
}

func (s *SQLiteDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s.DB = db
	return db, nil
}

func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLiteDatabase) CreateTable(model interface{}) error {
	return s.DB.AutoMigrate(model)
}

func (s *SQLiteDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := s.DB.Where(conditions).Find(model)
	return query, query.Error
}
