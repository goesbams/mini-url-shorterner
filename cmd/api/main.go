package main

import (
	"fmt"
	"log"
	"mini-url-shortener/config"
	"mini-url-shortener/internal/database"
	"mini-url-shortener/internal/routes"
	"mini-url-shortener/internal/server"
)

func main() {
	// load config.local.yaml file
	config.Init()
	cfg := config.GetConfig()

	// construct datasource name (DSN)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// initialize database connection
	db, err := database.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	fmt.Println(db)
	// TODO: setup services, handlers, and server

	// setup routes
	routes := routes.SetupRoutes()

	// start server
	serverPort := cfg.Server.Port
	srv := server.NewServer(routes, serverPort)

	log.Printf("Server started on port %s", serverPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
