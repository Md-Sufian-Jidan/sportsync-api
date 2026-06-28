package dto

import (
	"time"
	"sportsync-api/internal/domain/admin"
)

type Response struct {
	ID           uint               `json:"id"`
	UserID       uint               `json:"user_id"`
	ZoneID       uint               `json:"zone_id"`
	LicensePlate string             `json:"license_plate"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	User         *UserResponse      `json:"user,omitempty"`
	Zone         *admin.ParkingZone `json:"zone,omitempty"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type MyReservationResponse struct {
	ID           uint           `json:"id"`
	LicensePlate string         `json:"license_plate"`
	Status       string         `json:"status"`
	Zone         MyZoneResponse `json:"zone"`
	CreatedAt    time.Time      `json:"created_at"`
}

type MyZoneResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
