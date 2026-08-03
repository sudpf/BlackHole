package db

import (
	"BlackHole/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDatabase struct {
	logLevel string
	logFile  string
	logSize  string
	link     string
	DB       *gorm.DB
}

func (s *SQLiteDatabase) Connect(connectionString string) (*gorm.DB, error) {
	var la *logrusAdapter
	if s.logFile != "" {
		sqlLogger := logrus.New()
		output, err := logger.RotatingWriter(s.logFile, s.logSize)
		if err != nil {
			return nil, fmt.Errorf("initialize sqlite sql log: %w", err)
		}
		sqlLogger.SetOutput(output)
		sqlLogger.SetFormatter(&CustomFormatter{})
		if level, err := logrus.ParseLevel(s.logLevel); err == nil {
			sqlLogger.SetLevel(level)
		}
		la = NewLogrusAdapter(sqlLogger)
	}

	cfg := &gorm.Config{}
	if la != nil {
		cfg.Logger = la
	}

	db, err := gorm.Open(sqlite.Open(connectionString), cfg)
	if err != nil {
		return nil, err
	}
	s.DB = db
	return db, nil
}

func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLiteDatabase) CreateTable(model ...interface{}) error {
	return s.DB.AutoMigrate(model...)
}

func (s *SQLiteDatabase) CreateDatabase(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	dbExist, err := SQLiteDatabaseExist(s.link)
	if err != nil {
		return err
	}

	if dbExist {
		return nil
	}

	db, err := sql.Open("sqlite3", s.link)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

func (s *SQLiteDatabase) Query(ctx context.Context, model interface{}, conditions map[string]interface{}, options *QueryOptions) (*gorm.DB, error) {
	query := applyQueryOptions(s.DB.WithContext(ctx).Where(conditions), options).Find(model)
	return query, query.Error
}

func (s *SQLiteDatabase) QueryEx(ctx context.Context, model interface{}, conditions interface{}, options *QueryOptions) (*gorm.DB, error) {
	conditionMap, err := StructToConditions(conditions)
	if err != nil {
		return nil, err
	}

	query := applyQueryOptions(s.DB.WithContext(ctx).Where(conditionMap), options).Find(model)
	return query, query.Error
}

func (s *SQLiteDatabase) Insert(ctx context.Context, model interface{}) error {
	return s.DB.WithContext(ctx).Create(model).Error
}

func (s *SQLiteDatabase) Update(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return s.DB.WithContext(ctx).Model(model).Where(conditions).Updates(model).Error
}

func (s *SQLiteDatabase) Delete(ctx context.Context, model interface{}, conditions map[string]interface{}) error {
	return s.DB.WithContext(ctx).Where(conditions).Delete(model).Error
}

func NewSQLiteDatabase(connectionString string, logLevel string, logFile string, logSize string) (*SQLiteDatabase, error) {
	db := &SQLiteDatabase{logLevel: logLevel, logFile: logFile, logSize: logSize, link: connectionString}
	sqliteDb, err := db.Connect(connectionString)
	if err != nil {
		return nil, err
	}
	db.DB = sqliteDb
	return db, nil
}

func SQLiteDatabaseExist(connectionString string) (bool, error) {
	_, err := os.Stat(connectionString)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
