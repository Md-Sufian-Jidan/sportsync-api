package main

import (
	"sportsync-api/internal/config"
	"sportsync-api/internal/server"
)

func main() {
	// @title SpotSync API
// @version 1.0
// @description SpotSync Parking Management API
// @termsOfService http://swagger.io/terms/

// @contact.name SpotSync Team
// @contact.email support@spotsync.com

// @license.name MIT

// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
	// Load the config
	cfg := config.LoadEnv()
	// Connect to the database
	db := config.ConnectDatabase(cfg)
	// start the server
	server.Start(db, cfg)
}
