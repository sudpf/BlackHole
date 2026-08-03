package db

import (
	"BlackHole/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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

func (p *PostgreSQLDatabase) CreateDatabase(ctx context.Context) error {
	dbConfig, err := pgx.ParseConfig(p.link)
	if err != nil {
		return fmt.Errorf("parse postgres dsn: %w", err)
	}
	if dbConfig.Config.Database == "" {
		return fmt.Errorf("parse postgres dsn: database name is required")
	}

	dbName := dbConfig.Config.Database
	adminConfig := dbConfig.Copy()
	adminConfig.Config.Database = "postgres"

	db := stdlib.OpenDB(*adminConfig)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect postgres admin database: %w", err)
	}

	dbExist, err := postgresDatabaseExists(ctx, db, dbName)
	if err != nil {
		return fmt.Errorf("check postgres database %q: %w", dbName, err)
	}

	if dbExist {
		return nil
	}

	if _, err := db.ExecContext(ctx, "CREATE DATABASE "+quotePostgresIdentifier(dbName)); err != nil {
		return fmt.Errorf("create postgres database %q: %w", dbName, err)
	}
	return nil
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

func NewPostgreSQLDatabase(ctx context.Context, connectionString string, logLevel string, logFile string, logSize string) (*PostgreSQLDatabase, error) {
	db := &PostgreSQLDatabase{logLevel: logLevel, logFile: logFile, logSize: logSize, link: connectionString}

	if err := db.CreateDatabase(ctx); err != nil {
		return nil, err
	}

	pgDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = pgDb
	return db, nil
}

func postgresDatabaseExists(ctx context.Context, db *sql.DB, dbName string) (bool, error) {
	var value int
	err := db.QueryRowContext(ctx, "SELECT 1 FROM pg_database WHERE datname = $1", dbName).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func quotePostgresIdentifier(identifier string) string {
	return pgx.Identifier{identifier}.Sanitize()
}
