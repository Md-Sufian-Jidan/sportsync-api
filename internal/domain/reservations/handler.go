package reservations

import (
	"errors"
	"net/http"
	"sportsync-api/internal/domain/reservations/dto"
	"sportsync-api/internal/httpResponse"
	"strconv"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpResponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	reservation, err := h.service.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusBadRequest, httpResponse.Error{
				Success: false,
				Message: "Parking zone is full",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to create reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpResponse.Success{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    reservation,
	})
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpResponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to get reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
	})
}

func (h *handler) CancelReservation(c *echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  err.Error(),
		})
	}

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpResponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user id in context",
		})
	}
	userRole, _ := c.Get("user_role").(string)

	err = h.service.CancelReservation(uint(idParam), userID, userRole)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			return c.JSON(http.StatusForbidden, httpResponse.Error{
				Success: false,
				Message: "Forbidden",
				Errors:  err.Error(),
			})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpResponse.Error{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to cancel reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

func (h *handler) GetAllReservations(c *echo.Context) error {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to get reservations",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "reservations retrieved successfully",
		Data:    reservations,
	})
}

func (h *handler) GetReservationByID(c *echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.GetReservationByID(uint(idParam))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpResponse.Error{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to get reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation retrieved successfully",
		Data:    response,
	})
}

func (h *handler) UpdateReservation(c *echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid parking zone ID",
			Errors:  err.Error(),
		})
	}

	var req dto.UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	reservation, err := h.service.UpdateReservation(uint(idParam), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpResponse.Error{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to update reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "reservation updated successfully",
		Data:    reservation,
	})
}

func (h *handler) DeleteReservation(c *echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpResponse.Error{
			Success: false,
			Message: "Invalid reservation ID",
			Errors:  err.Error(),
		})
	}

	err = h.service.DeleteReservation(uint(idParam))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, httpResponse.Error{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpResponse.Error{
			Success: false,
			Message: "Failed to delete reservation",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation deleted successfully",
	})
}
