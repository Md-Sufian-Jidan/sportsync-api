package admin

import (
	"net/http"
	"sportsync-api/internal/domain/admin/dto"
	"sportsync-api/internal/httpResponse"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateParkingZone(ctx *echo.Context) error {
	var req dto.CreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	parkingZone, err := h.service.CreateParkingZone(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to create parking zone",
			Errors:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    parkingZone,
	})
}

func (h *handler) GetParkingZones(ctx *echo.Context) error {
	parkingZones, err := h.service.GetAllParkingZones()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to get parking zones",
			Errors:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    parkingZones,
	})
}

func (h *handler) GetParkingZoneByID(ctx *echo.Context) error {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.GetParkingZoneByID(uint(idParam))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to get parking zone",
			Errors:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    response,
	})
}

func (h *handler) UpdateParkingZone(ctx *echo.Context) error {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors:  err.Error(),
		})
	}

	var req dto.UpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	parkingZone, err := h.service.UpdateParkingZone(uint(idParam), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to update parking zone",
			Errors:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone updated successfully",
		Data:    parkingZone,
	})
}

func (h *handler) DeleteParkingZone(ctx *echo.Context) error {
	idParam, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors:  err.Error(),
		})
	}

	err = h.service.DeleteParkingZone(uint(idParam))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to delete parking zone",
			Errors:  err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}
