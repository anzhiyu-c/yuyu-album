package model

import (
	"time"

	"gorm.io/gorm"
)

type Setting struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ConfigKey string         `gorm:"column:config_key;type:varchar(100);unique;not null;comment:配置键" json:"key"`
	Value     string         `gorm:"type:text;comment:配置值" json:"value"`
	Comment   string         `gorm:"type:varchar(255);comment:配置注释" json:"comment"`
}
