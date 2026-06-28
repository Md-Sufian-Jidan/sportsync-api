package server

import (
	"fmt"
	"net/http"
	"sportsync-api/internal/config"
	"sportsync-api/internal/domain/admin"
	"sportsync-api/internal/domain/reservations"
	"sportsync-api/internal/domain/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()

	db.AutoMigrate(&user.User{}, &admin.ParkingZone{}, &reservations.Reservation{})
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "SportSync server is running successfully!")
	})

	user.RegisterRoutes(e, db, cfg)
	admin.RegisterRoutes(e, db, cfg)
	reservations.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Println("port", port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
