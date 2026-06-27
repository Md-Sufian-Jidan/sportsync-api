package admin

import (
	"sportsync-api/internal/auth"
	"sportsync-api/internal/domain/admin/dto"
)

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(r Repository, jwtService auth.JWTService) *service {
	return &service{repo: r, jwtService: jwtService}
}

func (s *service) CreateParkingZone(req dto.CreateRequest) (*dto.Response, error) {
	parkingZone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.repo.CreateParkingZone(&parkingZone); err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:            parkingZone.ID,
		Name:          parkingZone.Name,
		Type:          parkingZone.Type,
		TotalCapacity: parkingZone.TotalCapacity,
		PricePerHour:  parkingZone.PricePerHour,
		CreatedAt:     parkingZone.CreatedAt.String(),
		UpdatedAt:     parkingZone.UpdatedAt.String(),
	}
	return &response, nil
}

func (s *service) GetParkingZoneByID(id uint) (*dto.Response, error) {
	parkingZone, err := s.repo.GetParkingZoneByID(id)
	if err != nil {
		return nil, err
	}
	response := dto.Response{
		ID:            parkingZone.ID,
		Name:          parkingZone.Name,
		Type:          parkingZone.Type,
		TotalCapacity: parkingZone.TotalCapacity,
		PricePerHour:  parkingZone.PricePerHour,
		CreatedAt:     parkingZone.CreatedAt.String(),
		UpdatedAt:     parkingZone.UpdatedAt.String(),
	}
	return &response, nil
}

func (s *service) UpdateParkingZone(id uint, req dto.UpdateRequest) (*dto.Response, error) {
	parkingZone, err := s.repo.GetParkingZoneByID(id)
	if err != nil {
		return nil, err
	}
	parkingZone.Name = req.Name
	parkingZone.Type = req.Type
	parkingZone.TotalCapacity = req.TotalCapacity
	parkingZone.PricePerHour = req.PricePerHour

	err = s.repo.UpdateParkingZone(parkingZone)
	if err != nil {
		return nil, err
	}
	response := dto.Response{
		ID:            parkingZone.ID,
		Name:          parkingZone.Name,
		Type:          parkingZone.Type,
		TotalCapacity: parkingZone.TotalCapacity,
		PricePerHour:  parkingZone.PricePerHour,
		CreatedAt:     parkingZone.CreatedAt.String(),
		UpdatedAt:     parkingZone.UpdatedAt.String(),
	}
	return &response, nil
}

func (s *service) DeleteParkingZone(id uint) error {
	err := s.repo.DeleteParkingZone(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetAllParkingZones() ([]dto.Response, error) {
	parkingZones, err := s.repo.GetAllParkingZones()
	if err != nil {
		return nil, err
	}
	var responses []dto.Response
	for _, parkingZone := range parkingZones {
		responses = append(responses, dto.Response{
			ID:            parkingZone.ID,
			Name:          parkingZone.Name,
			Type:          parkingZone.Type,
			TotalCapacity: parkingZone.TotalCapacity,
			PricePerHour:  parkingZone.PricePerHour,
			CreatedAt:     parkingZone.CreatedAt.String(),
			UpdatedAt:     parkingZone.UpdatedAt.String(),
		})
	}
	return responses, nil
}
