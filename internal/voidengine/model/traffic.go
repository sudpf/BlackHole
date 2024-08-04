package model

import (
	"time"
)

type NetworkTraffic struct {
	ID              int64     `gorm:"primaryKey"`
	Timestamp       time.Time `gorm:"not null"`
	SourceIP        string    `gorm:"type:varchar;not null"`
	DestinationIP   string    `gorm:"type:varchar;not null"`
	SourcePort      int       `gorm:"not null"`
	DestinationPort int       `gorm:"not null"`
	Protocol        string    `gorm:"type:varchar;not null"`
	BytesIn         int64     `gorm:"not null"`
	BytesOut        int64     `gorm:"not null"`
	PacketCount     int       `gorm:"not null"`
	Description     string    `gorm:"type:varchar"`
}
