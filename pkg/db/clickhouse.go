package db

import (
	"database/sql"
	"fmt"
	"net/url"

	clickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseDatabase struct {
	logLevel string
	logFile  string
	link     string
	DB       *gorm.DB
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

func (c *ClickHouseDatabase) CreateTable(model ...interface{}) error {
	return c.DB.AutoMigrate(model...)
}

func (c *ClickHouseDatabase) CreateDatabase() error {
	u, err := url.Parse(c.link)
	if err != nil {
		return err
	}

	dbExist, err := ClickHouseDatabaseExist(fmt.Sprintf("%s:%d", u.Hostname(), u.Port), u.Query().Get("database"))
	if err != nil {
		return err
	}

	if dbExist {
		return nil
	}

	db, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s", fmt.Sprintf("%s:%d", u.Hostname(), u.Port)))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", u.Query().Get("database")))

	return err
}

func (c *ClickHouseDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := c.DB.Where(conditions).Find(model)
	return query, query.Error
}

func NewClickHouseDatabase(connectionString string, logLevel string, logFile string) (*ClickHouseDatabase, error) {
	db := &ClickHouseDatabase{logLevel: logLevel, logFile: logFile, link: connectionString}

	clickhouseDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = clickhouseDb
	return db, nil
}

func ClickHouseDatabaseExist(addr, database string) (bool, error) {
	db, err := sql.Open("clickhouse", fmt.Sprintf("tcp://%s", addr))
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT name FROM system.databases WHERE name = '%s'", database))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
