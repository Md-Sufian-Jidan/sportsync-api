package reservations

import (
	"sportsync-api/internal/domain/admin"
	"sportsync-api/internal/domain/user"

	"gorm.io/gorm"
)

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type Reservation struct {
	gorm.Model
	UserID       uint              `json:"user_id" gorm:"not null"`
	User         user.User         `gorm:"foreignKey:UserID"`
	ZoneID       uint              `json:"zone_id" gorm:"not null"`
	Zone         admin.ParkingZone `gorm:"foreignKey:ZoneID"`
	LicensePlate string            `json:"license_plate" gorm:"size:15;not null"`
	Status       string            `json:"status" gorm:"size:20;not null;default:active"`
}

