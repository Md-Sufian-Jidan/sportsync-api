package admin

import "gorm.io/gorm"

type ParkingZone struct {
	gorm.Model
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TotalCapacity int     `json:"total_capacity"`
	PricePerHour  float64 `json:"price_per_hour"`
}