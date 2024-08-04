package model

import (
	"BlackHole/pkg/config"
	"BlackHole/pkg/db"
)

var (
	cpDB *db.MySQLDatabase
	dpDB *db.ClickHouseDatabase
)

func InitDB(databaseConfig config.DatabaseConfig) error {
	if databaseConfig.MySQL != nil {
		mysqlDB, err := db.NewMySQLDatabase(databaseConfig.MySQL.Link, databaseConfig.MySQL.Debug, databaseConfig.MySQL.Log)
		if err != nil {
			panic(err)
		}
		if err := mysqlDB.CreateTable(&User{}); err != nil {
			panic(err)
		}
		cpDB = mysqlDB
	}

	if databaseConfig.ClickHouse != nil {
		ckDB, err := db.NewClickHouseDatabase(databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, databaseConfig.ClickHouse.Log)
		if err != nil {
			panic(err)
		}
		ckDB.CreateTable(&NetworkTraffic{})

		dpDB = ckDB
	}

	return nil
}

func ControlPlanDB() *db.MySQLDatabase {
	return cpDB
}

func DataPlanDB() *db.ClickHouseDatabase {
	return dpDB
}
