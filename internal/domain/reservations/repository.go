package reservations

import (
	"errors"
	"sportsync-api/internal/domain/admin"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is at full capacity")

type repository struct {
	db *gorm.DB
}

type Repository interface {
	CreateReservation(reservation *Reservation) error
	GetReservationByID(id uint) (*Reservation, error)
	GetReservationsByUserID(userID uint) ([]*Reservation, error)
	UpdateReservation(reservation *Reservation) error
	GetAllReservations() ([]*Reservation, error)
	DeleteReservation(id uint) error
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateReservation(reservation *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone admin.ParkingZone
		// 1. Lock the row!
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reservation.ZoneId).Error; err != nil {
			return err
		}

		// 2. Count current 'active' reservations for this zone
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", reservation.ZoneId, StatusActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check if active_count < zone.total_capacity
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create reservation
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		// Preload Zone and User on the created reservation so that the response fields are populated
		if err := tx.Preload("Zone").Preload("User").First(reservation, reservation.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) GetAllReservations() ([]*Reservation, error) {
	var reservations []*Reservation
	if err := r.db.Preload("Zone").Preload("User").Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) GetReservationByID(id uint) (*Reservation, error) {
	var reservation Reservation
	if err := r.db.Preload("Zone").Preload("User").First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) GetReservationsByUserID(userID uint) ([]*Reservation, error) {
	var reservations []*Reservation
	if err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) UpdateReservation(reservation *Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *repository) DeleteReservation(id uint) error {
	return r.db.Delete(&Reservation{}, id).Error
}
