package main

import (
	"fmt"
	"log"
	"mini-url-shortener/config"
	"mini-url-shortener/internal/database"
<<<<<<< HEAD
	"mini-url-shortener/internal/handlers"
	"mini-url-shortener/internal/repositories"
	"mini-url-shortener/internal/routes"
	"mini-url-shortener/internal/server"
	"mini-url-shortener/internal/services"
=======
	"mini-url-shortener/internal/routes"
	"mini-url-shortener/internal/server"
>>>>>>> e2e0deee502097425865ee87995fb2093e9e4e98
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

	// initialize repository, service, handler
	urlRepo := repositories.NewURLRepository(db)
	urlService := services.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService)

	// setup routes
	routes := routes.SetupRoutes(urlHandler)

	// start server
	serverPort := cfg.Server.Port
	srv := server.NewServer(routes, serverPort)

	log.Printf("Server started on port %s", serverPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
