package model

import (
	"BlackHole/pkg/config"
	"BlackHole/pkg/db"
	"errors"
	"fmt"
	"io"
)

type Models struct {
	User    *UserDAO
	Traffic *TrafficDAO
	closers []io.Closer
}

func New(databaseConfig config.DatabaseConfig, logDir string) (*Models, error) {
	models := &Models{}

	if databaseConfig.MySQL != nil {
		mysqlDB, err := db.NewMySQLDatabase(databaseConfig.MySQL.Link, databaseConfig.MySQL.Debug, logDir, databaseConfig.MySQL.Log)
		if err != nil {
			return nil, fmt.Errorf("initialize mysql: %w", err)
		}
		if err := mysqlDB.CreateTable(&User{}); err != nil {
			_ = mysqlDB.Close()
			return nil, fmt.Errorf("migrate user table: %w", err)
		}
		models.User = NewUserDAO(mysqlDB.DB)
		models.closers = append(models.closers, mysqlDB)
	}

	if databaseConfig.ClickHouse != nil {
		ckDB, err := db.NewClickHouseDatabase(databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, logDir, databaseConfig.ClickHouse.Log)
		if err != nil {
			_ = models.Close()
			return nil, fmt.Errorf("initialize clickhouse: %w", err)
		}
		if err := ckDB.CreateTable(&NetworkTraffic{}); err != nil {
			_ = ckDB.Close()
			_ = models.Close()
			return nil, fmt.Errorf("migrate network traffic table: %w", err)
		}
		models.Traffic = NewTrafficDAO(ckDB.DB)
		models.closers = append(models.closers, ckDB)
	}

	return models, nil
}

func (m *Models) Close() error {
	if m == nil {
		return nil
	}

	var err error
	for _, closer := range m.closers {
		if closeErr := closer.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	return err
}
