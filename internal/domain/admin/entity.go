package admin

import "gorm.io/gorm"

type ParkingZone struct {
	gorm.Model
	Name           string  `json:"name" gorm:"type:varchar(100);not null"`
	Type           string  `json:"type" gorm:"type:varchar(50);not null"`
	TotalCapacity  int     `json:"total_capacity" gorm:"not null"`
	PricePerHour   float64 `json:"price_per_hour" gorm:"type:decimal(10,2);not null"`
	AvailableSpots int     `json:"available_spots" gorm:"->"` // read-only field computed via subquery
}