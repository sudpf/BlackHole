package db

import (
	"BlackHole/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	mysqlParser "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDatabase struct {
	debug   bool
	logDir  string
	logFile string
	logSize string
	link    string
	DB      *gorm.DB
}

func (m *MySQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var la *logrusAdapter
	if m.debug {
		sqlLogger := logrus.New()
		output, err := logger.RotatingWriter(filepath.Join(m.logDir, m.logFile), m.logSize)
		if err != nil {
			return nil, fmt.Errorf("initialize mysql sql log: %w", err)
		}
		sqlLogger.SetOutput(output)
		sqlLogger.SetFormatter(&CustomFormatter{})
		la = NewLogrusAdapter(sqlLogger)
	}

	cfg := &gorm.Config{}
	if la != nil {
		cfg.Logger = la
	}

	db, err := gorm.Open(mysql.Open(connectionString), cfg)
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

func (m *MySQLDatabase) CreateDatabase(ctx context.Context) error {
	dbConfig, err := mysqlParser.ParseDSN(m.link)
	if err != nil {
		return fmt.Errorf("parse mysql dsn: %w", err)
	}
	if dbConfig.DBName == "" {
		return fmt.Errorf("parse mysql dsn: database name is required")
	}

	dbName := dbConfig.DBName
	dbConfig.DBName = ""
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("open mysql server connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect mysql server: %w", err)
	}

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", escapeMySQLIdentifier(dbName))); err != nil {
		return fmt.Errorf("create mysql database %q: %w", dbName, err)
	}
	return nil
}

func (m *MySQLDatabase) Query(ctx context.Context, model interface{}, conditions map[string]interface{}, options *QueryOptions) (*gorm.DB, error) {
	query := applyQueryOptions(m.DB.WithContext(ctx).Where(conditions), options).Find(model)
	return query, query.Error
}

func (m *MySQLDatabase) QueryEx(ctx context.Context, model interface{}, conditions interface{}, options *QueryOptions) (*gorm.DB, error) {
	query := applyQueryOptions(m.DB.WithContext(ctx).Where(conditions), options).Find(model)
	return query, query.Error
}

func (m *MySQLDatabase) Insert(ctx context.Context, model interface{}) error {
	return m.DB.WithContext(ctx).Create(model).Error
}

func (m *MySQLDatabase) Update(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return m.DB.WithContext(ctx).Model(model).Where(conditions).Updates(model).Error
}

func (m *MySQLDatabase) Delete(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return m.DB.WithContext(ctx).Where(conditions).Delete(model).Error
}

func NewMySQLDatabase(ctx context.Context, connectionString string, debug bool, logDir, logFile, logSize string) (*MySQLDatabase, error) {
	db := &MySQLDatabase{debug: debug, logDir: logDir, logFile: logFile, logSize: logSize, link: connectionString}
	if err := db.CreateDatabase(ctx); err != nil {
		log.Errorf("create mysql database error: %v", err)
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
