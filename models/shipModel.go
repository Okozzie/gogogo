package models

import (
	"gorm.io/gorm"
)

// Many to many relationship between Ship and Armaments, `ship_armaments` is the join table
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
