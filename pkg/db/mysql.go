package db

import (
	"BlackHole/pkg/config"
	"database/sql"
	"fmt"

	mysqlParser "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDatabase struct {
	debug   bool
	logFile string
	link    string
	DB      *gorm.DB
}

func (m *MySQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var logger logger.Interface
	if m.debug {
		logger := logrus.New()
		logger.SetOutput(&lumberjack.Logger{
			Filename: config.GetVoidEngineConfig().LogDir() + "/" + m.logFile,
			Compress: true,
		})
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, err
	}
	m.DB = db
	return db, nil
}

func (m *MySQLDatabase) Close() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (m *MySQLDatabase) CreateTable(model ...interface{}) error {
	return m.DB.AutoMigrate(model...)
}

func (m *MySQLDatabase) CreateDatabase() error {
	dbConfig, err := mysqlParser.ParseDSN(m.link)
	if err != nil {
		return err
	}

	log.Info(dbConfig)
	dbExist, err := MySQLDatabaseExist(dbConfig.Addr, dbConfig.User, dbConfig.Passwd, dbConfig.DBName)
	if err != nil {
		return err
	}

	if dbExist {
		return nil
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", dbConfig.User, dbConfig.Passwd, dbConfig.Addr))
	if err != nil {
		log.Info(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConfig.DBName))
	return err
}

func (m *MySQLDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := m.DB.Where(conditions).Find(model)
	return query, query.Error
}

func NewMySQLDatabase(connectionString string, debug bool, logFile string) (*MySQLDatabase, error) {
	db := &MySQLDatabase{debug: debug, logFile: logFile, link: connectionString}
	if err := db.CreateDatabase(); err != nil {
		log.Errorf("create database error:", err)
		return nil, err
	}

	mysqlDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = mysqlDb

	return db, nil
}

func MySQLDatabaseExist(addr, user, passwd, dbName string) (bool, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", user, passwd, addr))
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", dbName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
