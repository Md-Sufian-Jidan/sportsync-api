package reservations

import (
	"errors"
	"sportsync-api/internal/auth"
	"sportsync-api/internal/domain/reservations/dto"
)

var ErrForbidden = errors.New("forbidden: you can only cancel your own reservation")

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) *service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) CreateReservation(userID uint, req dto.CreateRequest) (*dto.Response, error) {
	res := Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       StatusActive,
	}
	err := s.repo.CreateReservation(&res)
	if err != nil {
		return nil, err
	}
	response := dto.Response{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}
	return &response, nil
}

func (s *service) GetAllReservations() ([]dto.Response, error) {
	reservations, err := s.repo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var responses []dto.Response
	for _, reservation := range reservations {
		res := dto.Response{
			ID:           reservation.ID,
			UserID:       reservation.UserID,
			ZoneID:       reservation.ZoneID,
			LicensePlate: reservation.LicensePlate,
			Status:       reservation.Status,
			CreatedAt:    reservation.CreatedAt,
			UpdatedAt:    reservation.UpdatedAt,
		}
		res.User = &dto.UserResponse{
			ID:    reservation.User.ID,
			Name:  reservation.User.Name,
			Email: reservation.User.Email,
			Role:  reservation.User.Role,
		}
		res.Zone = &dto.ZoneResponse{
			ID:            reservation.Zone.ID,
			Name:          reservation.Zone.Name,
			Type:          reservation.Zone.Type,
			TotalCapacity: reservation.Zone.TotalCapacity,
			PricePerHour:  reservation.Zone.PricePerHour,
		}
		responses = append(responses, res)
	}
	return responses, nil
}

func (s *service) GetReservationByID(id uint) (*dto.Response, error) {
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}
	response.User = &dto.UserResponse{
		ID:    reservation.User.ID,
		Name:  reservation.User.Name,
		Email: reservation.User.Email,
		Role:  reservation.User.Role,
	}
	response.Zone = &dto.ZoneResponse{
		ID:            reservation.Zone.ID,
		Name:          reservation.Zone.Name,
		Type:          reservation.Zone.Type,
		TotalCapacity: reservation.Zone.TotalCapacity,
		PricePerHour:  reservation.Zone.PricePerHour,
	}
	return &response, nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetReservationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.MyReservationResponse
	for _, res := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           res.ID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			Zone: dto.MyZoneResponse{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			CreatedAt: res.CreatedAt,
		})
	}
	return responses, nil
}

func (s *service) CancelReservation(id uint, userID uint, userRole string) error {
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return err
	}

	if userRole != "admin" && reservation.UserID != userID {
		return ErrForbidden
	}

	reservation.Status = StatusCancelled
	return s.repo.UpdateReservation(reservation)
}

func (s *service) UpdateReservation(id uint, req *dto.UpdateRequest) (*dto.Response, error) {
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	reservation.Status = req.Status
	err = s.repo.UpdateReservation(reservation)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}
	return &response, nil
}

func (s *service) DeleteReservation(id uint) error {
	err := s.repo.DeleteReservation(id)
	if err != nil {
		return err
	}
	return nil
}
