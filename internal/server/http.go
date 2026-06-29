package server

import (
	"errors"
	"fmt"
	"net/http"
	"sportsync-api/internal/config"
	"sportsync-api/internal/domain/admin"
	"sportsync-api/internal/domain/reservations"
	"sportsync-api/internal/domain/user"
	"sportsync-api/internal/httpResponse"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	_ "sportsync-api/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// ===============================
// @title SpotSync Server Health
// @description Health check endpoint
// @tags health
// @produce json
// @success 200 {string} string "OK"
// @router / [get]
// ===============================
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "SportSync server is running successfully!")
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal Server Error"
	var errorsDetail interface{} = err.Error()

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		if strMsg, ok := he.Message.(string); ok {
			message = strMsg
		} else {
			message = fmt.Sprintf("%v", he.Message)
		}
		errorsDetail = message
	} else if errors.Is(err, user.ErrorAlreadyExist) {
		code = http.StatusBadRequest
		message = "User already exists"
		errorsDetail = err.Error()
	} else if errors.Is(err, user.ErrInvalidCredentials) {
		code = http.StatusUnauthorized
		message = "Unauthorized"
		errorsDetail = err.Error()
	} else if errors.Is(err, reservations.ErrZoneFull) {
		code = http.StatusConflict
		message = "Parking zone is at full capacity"
		errorsDetail = err.Error()
	} else if errors.Is(err, reservations.ErrDuplicateLicensePlate) {
		code = http.StatusConflict
		message = "Duplicate license plate reservation"
		errorsDetail = err.Error()
	} else if errors.Is(err, reservations.ErrForbidden) {
		code = http.StatusForbidden
		message = "Forbidden"
		errorsDetail = err.Error()
	} else if errors.Is(err, admin.ErrParkingZoneNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
		code = http.StatusNotFound
		message = "Resource not found"
		errorsDetail = err.Error()
	} else {
		// Clean / mask internal/database errors to prevent leaking GORM errors to clients
		errorsDetail = "An unexpected error occurred"
	}

	_ = c.JSON(code, httpResponse.Error{
		Success: false,
		Message: message,
		Errors:  errorsDetail,
	})
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()

	db.AutoMigrate(&user.User{}, &admin.ParkingZone{}, &reservations.Reservation{})
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = CustomHTTPErrorHandler
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "SportSync server is running successfully!")
	})

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(302, "/swagger/index.html")
	})

	// Health endpoint (NOW SWAGGER WILL DETECT IT)
	e.GET("/health", healthCheck)

	user.RegisterRoutes(e, db, cfg)
	admin.RegisterRoutes(e, db, cfg)
	reservations.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Println("port", port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
