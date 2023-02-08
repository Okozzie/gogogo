package models

import "gorm.io/gorm"

type ShipArmament struct {
	gorm.Model `json:"-"`
	ShipID     uint
	ArmamentID uint
	Quantity   int `json:"quantity"`
}
