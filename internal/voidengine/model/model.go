package model

import (
	"BlackHole/internal/voidengine/config"
	"BlackHole/pkg/db"
	"context"
	"errors"
	"fmt"
	"io"
)

type Models struct {
	User    *UserDAO
	Traffic *TrafficDAO
	closers []io.Closer
}

func New(ctx context.Context, databaseConfig config.DatabaseConfig, logDir, logSize string) (*Models, error) {
	models := &Models{}

	if databaseConfig.MySQL != nil {
		mysqlDB, err := db.NewMySQLDatabase(ctx, databaseConfig.MySQL.Link, databaseConfig.MySQL.Debug, logDir, databaseConfig.MySQL.Log, logSize)
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
		ckDB, err := db.NewClickHouseDatabase(ctx, databaseConfig.ClickHouse.Link, databaseConfig.ClickHouse.Debug, logDir, databaseConfig.ClickHouse.Log, logSize)
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
