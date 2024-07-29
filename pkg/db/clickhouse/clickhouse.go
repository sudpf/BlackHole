package db

import (
	clickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseDatabase struct {
	DB *gorm.DB
}

func (c *ClickHouseDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(clickhouse.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	c.DB = db
	return db, nil
}

func (c *ClickHouseDatabase) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *ClickHouseDatabase) CreateTable(model interface{}) error {
	return c.DB.AutoMigrate(model)
}

func (c *ClickHouseDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := c.DB.Where(conditions).Find(model)
	return query, query.Error
}
