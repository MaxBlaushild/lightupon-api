package models

type Location struct {
	Latitude float64
	Longitude float64
	TripID uint
	Trip Trip
}

