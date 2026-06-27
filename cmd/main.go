package main

import (
	"sportsync-api/internal/config"
	"sportsync-api/internal/server"
)

func main() {
	// Load the config
	cfg := config.LoadEnv()
	// Connect to the database
	db := config.ConnectDatabase(cfg)
	// start the server
	server.Start(db, cfg)
}
