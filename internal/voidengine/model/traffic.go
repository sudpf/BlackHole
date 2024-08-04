package model

import (
	"time"
)

type NetworkTraffic struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"` // 唯一标识符
	Timestamp       time.Time `gorm:"not null"`                 // 记录时间
	SourceIP        string    `gorm:"size:45;not null"`         // 源IP地址（支持IPv4和IPv6）
	DestinationIP   string    `gorm:"size:45;not null"`         // 目的IP地址（支持IPv4和IPv6）
	SourcePort      int       `gorm:"not null"`                 // 源端口
	DestinationPort int       `gorm:"not null"`                 // 目的端口
	Protocol        string    `gorm:"size:10;not null"`         // 协议（如 TCP, UDP）
	BytesIn         int64     `gorm:"not null"`                 // 传入的字节数
	BytesOut        int64     `gorm:"not null"`                 // 传出的字节数
	PacketCount     int       `gorm:"not null"`                 // 包的数量
	Description     string    `gorm:"size:500"`                 // 其他描述信息
}
