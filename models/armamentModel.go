package models

import (
	"gorm.io/gorm"
)

type Armament struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"unique;not null" json:"name"`
	Quantity   int    `gorm:"-" json:"quantity"`
}
