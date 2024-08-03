package model

import (
	"BlackHole/pkg/config"
	"BlackHole/pkg/db"
)

func InitDB(databaseConfig config.DatabaseConfig) error {
	if databaseConfig.MySQL != nil {
		mysqlDB, err := db.NewMySQLDatabase(databaseConfig.MySQL.Link, databaseConfig.MySQL.Debug, databaseConfig.MySQL.Log)
		if err != nil {
			panic(err)
		}
		mysqlDB.CreateTable(&User{})
	}

	if databaseConfig.ClickHouse != nil {
		ckDB, err := db.NewClickHouseDatabase(databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, databaseConfig.ClickHouse.Log)
		if err != nil {
			panic(err)
		}
		ckDB.CreateTable(&User{})
	}

	return nil
}
