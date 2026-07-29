package model

import (
	"BlackHole/pkg/config"
	"BlackHole/pkg/db"
)

var (
	userDAO    *UserDAO
	trafficDAO *TrafficDAO
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
		userDAO = NewUserDAO(mysqlDB.DB)
	}

	if databaseConfig.ClickHouse != nil {
		ckDB, err := db.NewClickHouseDatabase(databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, databaseConfig.ClickHouse.Log)
		if err != nil {
			panic(err)
		}
		ckDB.CreateTable(&NetworkTraffic{})

		trafficDAO = NewTrafficDAO(ckDB.DB)
	}

	return nil
}

func GetUserDAO() *UserDAO {
	return userDAO
}

func GetTrafficDAO() *TrafficDAO {
	return trafficDAO
}
