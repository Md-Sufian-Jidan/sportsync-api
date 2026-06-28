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
//	@Summary		Reserve Parking Spot
//	@Description	Create a new parking reservation
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateRequest	true	"Reservation Information"
//	@Success		201		{object}	httpResponse.Success{data=dto.Response}
//	@Failure		400		{object}	httpResponse.Error
//	@Failure		401		{object}	httpResponse.Error
//	@Failure		409		{object}	httpResponse.Error
//	@Failure		500		{object}	httpResponse.Error
//	@Router			/reservations [post]

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

// GetMyReservations godoc
//
//	@Summary		Get My Reservations
//	@Description	Get all reservations of the authenticated user
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	httpResponse.Success{data=[]dto.Response}
//	@Failure		401	{object}	httpResponse.Error
//	@Failure		500	{object}	httpResponse.Error
//	@Router			/api/v1/reservations/my-reservations [get]

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

// CancelReservation godoc
//
//	@Summary		Cancel Reservation
//	@Description	Cancel an existing reservation. Drivers can only cancel their own reservations.
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Reservation ID"
//	@Success		200	{object}	httpResponse.Success
//	@Failure		400	{object}	httpResponse.Error
//	@Failure		401	{object}	httpResponse.Error
//	@Failure		403	{object}	httpResponse.Error
//	@Failure		404	{object}	httpResponse.Error
//	@Failure		500	{object}	httpResponse.Error
//	@Router			/api/v1/reservations/{id} [delete]

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

// GetAllReservations godoc
//
//	@Summary		Get All Reservations
//	@Description	Retrieve all reservations in the system (Admin only)
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	httpResponse.Success{data=[]dto.Response}
//	@Failure		401	{object}	httpResponse.Error
//	@Failure		403	{object}	httpResponse.Error
//	@Failure		500	{object}	httpResponse.Error
//	@Router			/api/v1/reservations [get]

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

// GetReservationByID godoc
//
//	@Summary		Get Reservation By ID
//	@Description	Get reservation details by reservation ID
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Reservation ID"
//	@Success		200	{object}	httpResponse.Success{data=dto.Response}
//	@Failure		400	{object}	httpResponse.Error
//	@Failure		401	{object}	httpResponse.Error
//	@Failure		404	{object}	httpResponse.Error
//	@Failure		500	{object}	httpResponse.Error
//	@Router			/api/v1/reservations/{id} [get]

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

// UpdateReservation godoc
//
//	@Summary		Update Reservation
//	@Description	Update an existing reservation
//	@Tags			Reservations
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Reservation ID"
//	@Param			request	body		dto.UpdateRequest	true	"Updated Reservation"
//	@Success		200		{object}	httpResponse.Success{data=dto.Response}
//	@Failure		400		{object}	httpResponse.Error
//	@Failure		401		{object}	httpResponse.Error
//	@Failure		404		{object}	httpResponse.Error
//	@Failure		500		{object}	httpResponse.Error
//	@Router			/api/v1/reservations/{id} [put]

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
