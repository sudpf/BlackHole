package model

import (
	"errors"

	"gorm.io/gorm"
)

type UserQuery struct {
	PageNo   int
	PageSize int
	OrderBy  string
	Username *string
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) List(query UserQuery) ([]User, error) {
	database := d.db
	if query.Username != nil {
		database = database.Where("name = ?", *query.Username)
	}
	if query.PageNo > 0 && query.PageSize > 0 {
		database = database.Offset((query.PageNo - 1) * query.PageSize).Limit(query.PageSize)
	}
	if query.OrderBy != "" {
		if query.OrderBy == "desc" {
			database = database.Order("id DESC")
		} else {
			database = database.Order("id ASC")
		}
	}

	var users []User
	if err := database.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDAO) Create(user *User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) FindByName(username string) (*User, error) {
	var user User
	err := d.db.Where("name = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Update(username string, user *User) error {
	return d.db.Model(user).Where("name = ?", username).Updates(user).Error
}

func (d *UserDAO) DeleteByName(username string) error {
	return d.db.Where("name = ?", username).Delete(&User{}).Error
}
