package db

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDatabase struct {
	logLevel string
	logFile  string
	link     string
	DB       *gorm.DB
}

func (p *PostgreSQLDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
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

func (p *PostgreSQLDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := p.DB.Where(conditions).Find(model)
	return query, query.Error
}

func (p *PostgreSQLDatabase) Insert(model interface{}) error {
	return p.DB.Create(model).Error
}

func (p *PostgreSQLDatabase) Update(model interface{}, conditions map[string]interface{}) error {
	return p.DB.Model(model).Where(conditions).Updates(model).Error
}

func (p *PostgreSQLDatabase) Delete(model interface{}, conditions map[string]interface{}) error {
	return p.DB.Where(conditions).Delete(model).Error
}

func NewPostgreSQLDatabase(connectionString string, logLevel string, logFile string) (*PostgreSQLDatabase, error) {
	db := &PostgreSQLDatabase{logLevel: logLevel, logFile: logFile, link: connectionString}

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
