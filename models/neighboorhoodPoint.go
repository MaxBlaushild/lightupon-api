package models

import(
    "github.com/jinzhu/gorm"
)

type NeighborhoodPoint struct {
	gorm.Model
	NeighborhoodID uint
	Latitude float64
	Longitude float64
}

func getNeighborhoodIDForLocation (lat string, lon string) uint {
	neighborhoodPoint := NeighborhoodPoint{}
	DB.Order("((latitude - " + lat + ")^2.0 + ((longitude - " + lon + ")* cos(latitude / 57.3))^2.0) asc").First(&neighborhoodPoint)
	return neighborhoodPoint.NeighborhoodID
}
