package admin

import (
	"sportsync-api/internal/auth"
	"sportsync-api/internal/config"
	"sportsync-api/internal/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	adminRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JwtSecret)
	adminService := NewService(adminRepository, jwtService)
	adminHandler := NewHandler(adminService)

	api := e.Group("/api/v1")

	api.POST("/zones", adminHandler.CreateParkingZone, middleware.CheckAuth(jwtService))
	api.GET("/zones", adminHandler.GetParkingZones)
	api.GET("/zones/:id", adminHandler.GetParkingZoneByID)
	api.PUT("/zones/:id", adminHandler.UpdateParkingZone, middleware.CheckAuth(jwtService))
	api.DELETE("/zones/:id", adminHandler.DeleteParkingZone, middleware.CheckAuth(jwtService))
}
