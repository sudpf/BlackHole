package model

import (
	"BlackHole/pkg/config"
	"BlackHole/pkg/db"
	"fmt"
)

type Models struct {
	User    *UserDAO
	Traffic *TrafficDAO
}

func New(databaseConfig config.DatabaseConfig) (*Models, error) {
	models := &Models{}

	if databaseConfig.MySQL != nil {
		mysqlDB, err := db.NewMySQLDatabase(databaseConfig.MySQL.Link, databaseConfig.MySQL.Debug, databaseConfig.MySQL.Log)
		if err != nil {
			return nil, fmt.Errorf("initialize mysql: %w", err)
		}
		if err := mysqlDB.CreateTable(&User{}); err != nil {
			return nil, fmt.Errorf("migrate user table: %w", err)
		}
		models.User = NewUserDAO(mysqlDB.DB)
	}

	if databaseConfig.ClickHouse != nil {
		ckDB, err := db.NewClickHouseDatabase(databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, databaseConfig.ClickHouse.Log)
		if err != nil {
			return nil, fmt.Errorf("initialize clickhouse: %w", err)
		}
		if err := ckDB.CreateTable(&NetworkTraffic{}); err != nil {
			return nil, fmt.Errorf("migrate network traffic table: %w", err)
		}
		models.Traffic = NewTrafficDAO(ckDB.DB)
	}

	return models, nil
}
