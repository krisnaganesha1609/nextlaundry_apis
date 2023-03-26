package models

import (
	"time"

	"gorm.io/gorm"
)

type Types string

const (
	User    Types = "users"
	Member  Types = "member"
	Outlet  Types = "outlet"
	Package Types = "paket"
	Login   Types = "login"
)

type LogHistory struct {
	IDLog     int            `gorm:"primaryKey;autoIncrement;column:id" json:"id_log"`
	Log       string         `gorm:"type:VARCHAR(250);column:log" json:"log_history"`
	Type      Types          `gorm:"type:VARCHAR(10);column:type" json:"log_type"`
	CreatedAt time.Time      `gorm:"type:DATETIME;column:created_at" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index,column:deleted_at" json:"deleted_at"`
}

func (l *LogHistory) TableName() string {
	return "log_history"
}
