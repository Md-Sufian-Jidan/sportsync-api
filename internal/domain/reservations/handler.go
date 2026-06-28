package reservations

import (
	"net/http"
	"sportsync-api/internal/domain/reservations/dto"
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

// CreateReservation godoc
//
// @Summary Reserve Parking Spot
// @Description Reserve a parking spot
// @Tags Reservations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateRequest true "Reservation"
// @Success 201 {object} httpResponse.Success
// @Failure 400 {object} httpResponse.Error
// @Failure 409 {object} httpResponse.Error
// @Router /reservations [post]

func (h *handler) CreateReservation(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user id in context")
	}

	var req dto.CreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	reservation, err := h.service.CreateReservation(userID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, httpResponse.Success{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    reservation,
	})
}

func (h *handler) GetMyReservations(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user id in context")
	}

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
	})
}

func (h *handler) CancelReservation(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user id in context")
	}
	userRole, _ := c.Get("user_role").(string)

	err = h.service.CancelReservation(uint(idParam), userID, userRole)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

func (h *handler) GetAllReservations(c echo.Context) error {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "reservations retrieved successfully",
		Data:    reservations,
	})
}

func (h *handler) GetReservationByID(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := h.service.GetReservationByID(uint(idParam))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation retrieved successfully",
		Data:    response,
	})
}

func (h *handler) UpdateReservation(c echo.Context) error {
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

	reservation, err := h.service.UpdateReservation(uint(idParam), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "reservation updated successfully",
		Data:    reservation,
	})
}

func (h *handler) DeleteReservation(c echo.Context) error {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.service.DeleteReservation(uint(idParam))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Reservation deleted successfully",
	})
}
