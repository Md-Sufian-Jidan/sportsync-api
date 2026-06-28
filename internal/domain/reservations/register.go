package reservations

import (
	"sportsync-api/internal/auth"
	"sportsync-api/internal/config"
	"sportsync-api/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	reservationRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JwtSecret)
	reservationService := NewService(reservationRepository, jwtService)
	reservationHandler := NewHandler(reservationService)

	api := e.Group("/api/v1")

	api.POST(
		"/reservations",
		reservationHandler.CreateReservation,
		middleware.CheckAuth(jwtService),
	)

	api.GET(
		"/reservations/my-reservations",
		reservationHandler.GetMyReservations,
		middleware.CheckAuth(jwtService),
	)

	api.GET(
		"/reservations",
		reservationHandler.GetAllReservations,
		middleware.CheckAuth(jwtService),
		middleware.RequireRole("admin"),
	)

	api.DELETE(
		"/reservations/:id",
		reservationHandler.CancelReservation,
		middleware.CheckAuth(jwtService),
	)
}
