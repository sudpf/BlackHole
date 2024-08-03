package db

import (
	"BlackHole/pkg/config"
	"fmt"

	clickhouseParser "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	clickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ClickHouseDatabase struct {
	debug   bool
	logFile string
	link    string
	DB      *gorm.DB
}

func (c *ClickHouseDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(clickhouse.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Info(err)
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
	var logger logger.Interface
	if c.debug {
		logger := logrus.New()
		logger.SetOutput(&lumberjack.Logger{
			Filename: config.GetVoidEngineConfig().LogDir() + "/" + c.logFile,
			Compress: true,
		})
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// 解析 DSN
	connParams, err := clickhouseParser.ParseDSN(c.link)
	if err != nil {
		log.Fatalf("failed to parse DSN: %v", err)
	}

	dsn := fmt.Sprintf("tcp://%s?debug=true&username=%s&password=%s&read_timeout=10s",
		connParams.Addr[0], connParams.Auth.Username, connParams.Auth.Password)
	log.Info(dsn)
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		log.Fatalf("failed to connect to ClickHouse: %v", err)
	}

	// 使用原生 SQL 创建新数据库
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", connParams.Auth.Database)).Error; err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	return err
}

func (c *ClickHouseDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := c.DB.Where(conditions).Find(model)
	return query, query.Error
}

func NewClickHouseDatabase(connectionString string, debug bool, logFile string) (*ClickHouseDatabase, error) {
	db := &ClickHouseDatabase{debug: debug, logFile: logFile, link: connectionString}

	if err := db.CreateDatabase(); err != nil {
		return nil, err
	}

	clickhouseDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = clickhouseDb
	return db, nil
}