package dto

type CreateRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type UpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active completed cancelled"`
}

