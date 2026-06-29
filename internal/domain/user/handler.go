package user

import (
	"net/http"
	"sportsync-api/internal/domain/user/dto"
	"sportsync-api/internal/httpResponse"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

// Register godoc
//
// @Summary Register User
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.CreateRequest true "Register Request"
// @Success 201 {object} httpResponse.Success
// @Failure 400 {object} httpResponse.Error
// @Failure 409 {object} httpResponse.Error
// @Router /auth/register [post]
func (h *handler) CreateUser(c echo.Context) error {
	var req dto.CreateRequest // input

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, httpResponse.Success{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

// LoginUser godoc
//
// @Summary Login User
// @Description Login with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} httpResponse.Success
// @Failure 401 {object} httpResponse.Error
// @Router /auth/login [post]
func (h *handler) LoginUser(c echo.Context) error {
	var req dto.LoginRequest // input

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// GetMe godoc
//
//	@Summary		Get authenticated user
//	@Description	Returns the currently authenticated user's information
//	@Tags			Authentication
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	httpResponse.Success{data=dto.Response}
//	@Failure		401	{object}	httpResponse.Error
//	@Router			/auth/me [get]
func (h *handler) GetMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user id in context")
	}
	email, _ := c.Get("user_email").(string)
	name, _ := c.Get("user_name").(string)
	role, _ := c.Get("user_role").(string)
	return c.JSON(http.StatusOK, httpResponse.Success{
		Success: true,
		Message: "User information fetched successfully",
		Data: dto.Response{
			ID:    userID,
			Name:  name,
			Email: email,
			Role:  role,
		},
	})
}
