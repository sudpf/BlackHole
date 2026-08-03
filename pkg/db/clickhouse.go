package db

import (
	"BlackHole/pkg/logger"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	clickhouseParser "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	clickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseDatabase struct {
	debug   bool
	logDir  string
	logFile string
	logSize string
	link    string
	DB      *gorm.DB
}

func (c *ClickHouseDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var la *logrusAdapter
	if c.debug {
		sqlLogger := logrus.New()
		output, err := logger.RotatingWriter(filepath.Join(c.logDir, c.logFile), c.logSize)
		if err != nil {
			return nil, fmt.Errorf("initialize clickhouse sql log: %w", err)
		}
		sqlLogger.SetOutput(output)
		sqlLogger.SetFormatter(&CustomFormatter{})
		sqlLogger.SetLevel(logrus.DebugLevel)
		la = NewLogrusAdapter(sqlLogger)
	}

	cfg := &gorm.Config{}
	if la != nil {
		cfg.Logger = la
	}

	db, err := gorm.Open(clickhouse.Open(connectionString), cfg)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	c.DB = db
	return db, nil
}

func (c *ClickHouseDatabase) Close() error {
	if c.DB == nil {
		return nil
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *ClickHouseDatabase) CreateTable(model ...interface{}) error {
	return c.DB.AutoMigrate(model...)
}

func (c *ClickHouseDatabase) CreateDatabase(ctx context.Context) error {
	var la *logrusAdapter
	if c.debug {
		sqlLogger := logrus.New()
		output, err := logger.RotatingWriter(filepath.Join(c.logDir, c.logFile), c.logSize)
		if err != nil {
			return fmt.Errorf("initialize clickhouse sql log: %w", err)
		}
		sqlLogger.SetOutput(output)
		sqlLogger.SetFormatter(&CustomFormatter{})
		sqlLogger.SetLevel(logrus.DebugLevel)
		la = NewLogrusAdapter(sqlLogger)
	}

	connParams, err := clickhouseParser.ParseDSN(c.link)
	if err != nil {
		return fmt.Errorf("parse clickhouse dsn: %w", err)
	}

	if len(connParams.Addr) == 0 {
		return fmt.Errorf("parse clickhouse dsn: missing address")
	}
	if connParams.Auth.Database == "" {
		return fmt.Errorf("parse clickhouse dsn: database name is required")
	}

	dsn := fmt.Sprintf("tcp://%s?username=%s&password=%s&read_timeout=10s",
		connParams.Addr[0], connParams.Auth.Username, connParams.Auth.Password)
	cfg := &gorm.Config{}
	if la != nil {
		cfg.Logger = la
	}

	db, err := gorm.Open(clickhouse.Open(dsn), cfg)
	if err != nil {
		return fmt.Errorf("connect clickhouse admin: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get clickhouse sql db: %w", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("connect clickhouse admin: %w", err)
	}

	if err := db.WithContext(ctx).Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", escapeClickHouseIdentifier(connParams.Auth.Database))).Error; err != nil {
		return fmt.Errorf("create clickhouse database %q: %w", connParams.Auth.Database, err)
	}

	return nil
}

func (c *ClickHouseDatabase) Query(ctx context.Context, model interface{}, conditions map[string]interface{}, options *QueryOptions) (*gorm.DB, error) {
	query := applyQueryOptions(c.DB.WithContext(ctx).Where(conditions), options).Find(model)
	return query, query.Error
}

func (c *ClickHouseDatabase) QueryEx(ctx context.Context, model interface{}, conditions interface{}, options *QueryOptions) (*gorm.DB, error) {
	conditionMap, err := StructToConditions(conditions)
	if err != nil {
		return nil, err
	}

	query := applyQueryOptions(c.DB.WithContext(ctx).Where(conditionMap), options).Find(model)
	return query, query.Error
}

func (c *ClickHouseDatabase) Insert(ctx context.Context, model interface{}) error {
	return c.DB.WithContext(ctx).Create(model).Error
}

func (c *ClickHouseDatabase) Update(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return c.DB.WithContext(ctx).Model(model).Where(conditions).Updates(model).Error
}

func (c *ClickHouseDatabase) Delete(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return c.DB.WithContext(ctx).Where(conditions).Delete(model).Error
}

func NewClickHouseDatabase(ctx context.Context, connectionString string, debug bool, logDir, logFile, logSize string) (*ClickHouseDatabase, error) {
	db := &ClickHouseDatabase{debug: debug, logDir: logDir, logFile: logFile, logSize: logSize, link: connectionString}

	if err := db.CreateDatabase(ctx); err != nil {
		return nil, err
	}

	clickhouseDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = clickhouseDb
	return db, nil
}

func escapeClickHouseIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "`", "``")
}
