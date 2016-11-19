package models

import(	"github.com/jinzhu/gorm")

type Location struct {
	gorm.Model
	Latitude float64
	Longitude float64
	TripID uint
}

