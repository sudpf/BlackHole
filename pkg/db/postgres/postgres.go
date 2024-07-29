package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDatabase struct {
	DB *gorm.DB
}

func (p *PostgreSQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	p.DB = db
	return db, nil
}

func (p *PostgreSQLDatabase) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgreSQLDatabase) CreateTable(model interface{}) error {
	return p.DB.AutoMigrate(model)
}

func (p *PostgreSQLDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := p.DB.Where(conditions).Find(model)
	return query, query.Error
}
