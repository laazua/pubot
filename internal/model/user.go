package model

import (
	"time"

	"gorm.io/gorm"
)

type PbUser struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(255);not null"`
	Password  string `gorm:"type:varchar(512);not null"`
	Role      string `gorm:"type:varchar(64);not null"` // admin or user
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (PbUser) TableName() string {
	return "pb_user"
}
