package db

import (
	"BlackHole/pkg/logger"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDatabase struct {
	logLevel string
	logFile  string
	logSize  string
	link     string
	DB       *gorm.DB
}

func (p *PostgreSQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var la *logrusAdapter
	if p.logFile != "" {
		sqlLogger := logrus.New()
		output, err := logger.RotatingWriter(p.logFile, p.logSize)
		if err != nil {
			return nil, fmt.Errorf("initialize postgres sql log: %w", err)
		}
		sqlLogger.SetOutput(output)
		sqlLogger.SetFormatter(&CustomFormatter{})
		if level, err := logrus.ParseLevel(p.logLevel); err == nil {
			sqlLogger.SetLevel(level)
		}
		la = NewLogrusAdapter(sqlLogger)
	}

	cfg := &gorm.Config{}
	if la != nil {
		cfg.Logger = la
	}

	db, err := gorm.Open(postgres.Open(connectionString), cfg)
	if err != nil {
		return nil, err
	}
	p.DB = db
	return db, nil
}

func (p *PostgreSQLDatabase) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *PostgreSQLDatabase) CreateTable(model ...interface{}) error {
	return p.DB.AutoMigrate(model...)
}

func (p *PostgreSQLDatabase) CreateDatabase() error {
	dbConfig, err := pgx.ParseConfig(p.link)
	if err != nil {
		return err
	}

	dbExist, err := PGDatabaseExist(dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database)
	if err != nil {
		return err
	}

	if dbExist {
		return nil
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database))
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConfig.Database))
	return err
}

func (p *PostgreSQLDatabase) Query(ctx context.Context, model interface{}, conditions map[string]interface{}, options *QueryOptions) (*gorm.DB, error) {
	query := applyQueryOptions(p.DB.WithContext(ctx).Where(conditions), options).Find(model)
	return query, query.Error
}

func (p *PostgreSQLDatabase) QueryEx(ctx context.Context, model interface{}, conditions interface{}, options *QueryOptions) (*gorm.DB, error) {
	conditionMap, err := StructToConditions(conditions)
	if err != nil {
		return nil, err
	}

	query := applyQueryOptions(p.DB.WithContext(ctx).Where(conditionMap), options).Find(model)
	return query, query.Error
}

func (p *PostgreSQLDatabase) Insert(ctx context.Context, model interface{}) error {
	return p.DB.WithContext(ctx).Create(model).Error
}

func (p *PostgreSQLDatabase) Update(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return p.DB.WithContext(ctx).Model(model).Where(conditions).Updates(model).Error
}

func (p *PostgreSQLDatabase) Delete(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return p.DB.WithContext(ctx).Where(conditions).Delete(model).Error
}

func NewPostgreSQLDatabase(connectionString string, logLevel string, logFile string, logSize string) (*PostgreSQLDatabase, error) {
	db := &PostgreSQLDatabase{logLevel: logLevel, logFile: logFile, logSize: logSize, link: connectionString}

	pgDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = pgDb
	return db, nil
}

func PGDatabaseExist(addr, user, passwd, dbName string) (bool, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=postgres", addr, user, passwd))
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT datname FROM pg_database WHERE datname = '%s'", dbName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}
