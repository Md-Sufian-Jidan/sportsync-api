package main

import (
	"sportsync-api/internal/config"
	"sportsync-api/internal/server"
)

// @title SpotSync API
// @version 1.0
// @description SpotSync Parking Management API
// @termsOfService http://swagger.io/terms/

// @contact.name SpotSync Team
// @contact.email support@spotsync.com

// @license.name MIT

// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)
	server.Start(db, cfg)
}
