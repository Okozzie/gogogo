package models

import (
	"gorm.io/gorm"
)

// User has and belongs to many languages, `user_languages` is the join table
type Ship struct {
	gorm.Model
	Name      string     `json:"name"`
	Class     string     `json:"class"`
	Crew      uint       `json:"crew"`
	Image     string     `json:"image"`
	Value     float32    `json:"value"`
	Status    string     `json:"status"`
	Armaments []Armament `gorm:"many2many:ship_armaments;" json:"armaments"`
}
