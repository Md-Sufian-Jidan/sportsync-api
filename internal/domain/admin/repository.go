package admin

import (
	"errors"

	"gorm.io/gorm"
)

var ErrParkingZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	CreateParkingZone(parkingZone *ParkingZone) error
	GetParkingZoneByID(id uint) (*ParkingZone, error)
	UpdateParkingZone(parkingZone *ParkingZone) error
	GetAllParkingZones() ([]*ParkingZone, error)
	DeleteParkingZone(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r repository) CreateParkingZone(parkingZone *ParkingZone) error {
	return r.db.Create(parkingZone).Error
}

func (r repository) GetAllParkingZones() ([]*ParkingZone, error) {
	var parkingZones []*ParkingZone
	err := r.db.Find(&parkingZones).Error
	if err != nil {
		return nil, err
	}
	return parkingZones, nil
}

func (r repository) GetParkingZoneByID(id uint) (*ParkingZone, error) {
	var parkingZone ParkingZone
	err := r.db.First(&parkingZone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParkingZoneNotFound
		}
		return nil, err
	}
	return &parkingZone, nil
}

func (r repository) UpdateParkingZone(parkingZone *ParkingZone) error {
	return r.db.Save(parkingZone).Error
}

func (r repository) DeleteParkingZone(id uint) error {
	return r.db.Delete(&ParkingZone{}, id).Error
}
