package model

import (
	"context"

	"gorm.io/gorm"
)

type TrafficQuery struct {
	PageNo   int
	PageSize int
	OrderBy  string
}

type TrafficDAO struct {
	db *gorm.DB
}

func NewTrafficDAO(db *gorm.DB) *TrafficDAO {
	return &TrafficDAO{db: db}
}

func (d *TrafficDAO) List(ctx context.Context, query TrafficQuery) ([]NetworkTraffic, error) {
	database := d.db.WithContext(ctx)
	if query.PageNo > 0 && query.PageSize > 0 {
		database = database.Offset((query.PageNo - 1) * query.PageSize).Limit(query.PageSize)
	}
	if query.OrderBy != "" {
		if query.OrderBy == "desc" {
			database = database.Order("timestamp DESC")
		} else {
			database = database.Order("timestamp ASC")
		}
	}

	var traffics []NetworkTraffic
	if err := database.Find(&traffics).Error; err != nil {
		return nil, err
	}
	return traffics, nil
}
