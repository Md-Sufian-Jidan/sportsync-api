package admin

import (
	"net/http"
	"sportsync-api/internal/domain/admin/dto"
	"sportsync-api/internal/httpResponse"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

// CreateParkingZone godoc
//
// @Summary Create Parking Zone
// @Description Create a new parking zone
// @Tags Parking Zones
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateRequest true "Parking Zone"
// @Success 201 {object} httpResponse.Success
// @Failure 400 {object} httpResponse.Error
// @Failure 401 {object} httpResponse.Error
// @Router /zones [post]

func (h *handler) CreateParkingZone(c echo.Context) error {
	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	parkingZone, err := h.service.CreateParkingZone(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, httpResponse.Success{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    parkingZone,
	})
}

func (h *handler) GetParkingZones(c echo.Context) error {
	parkingZones, err := h.service.GetAllParkingZones()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    parkingZones,
	})
}

func (h *handler) GetParkingZoneByID(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := h.service.GetParkingZoneByID(uint(idParam))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    response,
	})
}

func (h *handler) UpdateParkingZone(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	parkingZone, err := h.service.UpdateParkingZone(uint(idParam), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone updated successfully",
		Data:    parkingZone,
	})
}

func (h *handler) DeleteParkingZone(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.service.DeleteParkingZone(uint(idParam))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}
