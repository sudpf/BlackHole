package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	mysqlParser "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
	debug   bool
	logDir  string
	logFile string
	link    string
	DB      *gorm.DB
}

func (m *MySQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var la *logrusAdapter
	if m.debug {
		logger := logrus.New()
		logger.SetOutput(&lumberjack.Logger{
			Filename: filepath.Join(m.logDir, m.logFile),
			Compress: true,
		})
		logger.SetFormatter(&CustomFormatter{})
		la = NewLogrusAdapter(logger)
	}

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: la})
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

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", dbConfig.User, dbConfig.Passwd, dbConfig.Addr))
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", escapeMySQLIdentifier(dbConfig.DBName)))
	return err
}

func (m *MySQLDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	pageNo, okPageNo := conditions["PageNo"].(int)
	pageSize, okPageSize := conditions["PageSize"].(int)
	order, okOrder := conditions["OrderBy"].(string)

	delete(conditions, "PageNo")
	delete(conditions, "PageSize")
	delete(conditions, "OrderBy")

	db := m.DB.Where(conditions)
	if okPageNo && okPageSize {
		db = db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	}

	if okOrder {
		if order == "desc" {
			db = db.Order("id DESC")
		} else {
			db = db.Order("id ASC")
		}
	}

	db = db.Find(model)

	return db, db.Error
}

func (m *MySQLDatabase) QueryEx(model interface{}, conditions interface{}) (*gorm.DB, error) {
	query := m.DB.Where(conditions).Find(model)
	return query, query.Error
}

func (m *MySQLDatabase) Insert(model interface{}) error {
	return m.DB.Create(model).Error
}

func (m *MySQLDatabase) Update(model interface{}, conditions map[string]interface{}) error {
	return m.DB.Model(model).Where(conditions).Updates(model).Error
}

func (m *MySQLDatabase) Delete(model interface{}, conditions map[string]interface{}) error {
	return m.DB.Where(conditions).Delete(model).Error
}

func NewMySQLDatabase(connectionString string, debug bool, logDir, logFile string) (*MySQLDatabase, error) {
	db := &MySQLDatabase{debug: debug, logDir: logDir, logFile: logFile, link: connectionString}
	if err := db.CreateDatabase(); err != nil {
		log.Errorf("create database error: %v", err)
		return nil, err
	}

	mysqlDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = mysqlDb

	return db, nil
}

func escapeMySQLIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "`", "``")
}
