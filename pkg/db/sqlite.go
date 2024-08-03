package db

import (
	"database/sql"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDatabase struct {
	logLevel string
	logFile  string
	link     string
	DB       *gorm.DB
}

func (s *SQLiteDatabase) Connect(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
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

func (s *SQLiteDatabase) CreateDatabase() error {
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

func (s *SQLiteDatabase) Query(model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	query := s.DB.Where(conditions).Find(model)
	return query, query.Error
}

func (s *SQLiteDatabase) Insert(model interface{}) error {
	return s.DB.Create(model).Error
}

func (s *SQLiteDatabase) Update(model interface{}, conditions map[string]interface{}) error {
	return s.DB.Model(model).Where(conditions).Updates(model).Error
}

func (s *SQLiteDatabase) Delete(model interface{}, conditions map[string]interface{}) error {
	return s.DB.Where(conditions).Delete(model).Error
}

func NewSQLiteDatabase(connectionString string, logLevel string, logFile string) (*SQLiteDatabase, error) {
	db := &SQLiteDatabase{logLevel: logLevel, logFile: logFile, link: connectionString}
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
