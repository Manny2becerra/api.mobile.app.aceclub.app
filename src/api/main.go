package main

import (
	"api-mobile-app/src/api/common/config"
	"api-mobile-app/src/api/common/database"
	logging "api-mobile-app/src/api/common/logging"
	establishments "api-mobile-app/src/api/establishments"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")
	// we need to create a new http server
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbClient, err := database.Connect(config.DB.ConnectionString, config.DB.DbName, config.EstablishmentProfileImageBucket, config.AWSRegion, config.AWSAccessKey, config.AWSSecretKey)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	logger := logging.NewLogger()
	logger.StartUpAPI()
	router := http.NewServeMux()

	establishmentsRouter := establishments.EstablishmentsRouter()
	router.Handle("POST /establishments/", http.StripPrefix("/establishments", establishmentsRouter))

	httpServer := &http.Server{
		Addr:    config.Port,
		Handler: router,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
	// we need to initialize listening on port 3005
	// we need to create middleware for the server which authenticates the request
	// we need to create a server multiplexer for the main routes
	// we need to load all environment configurations
	// we need to create an establishment multiplexer for the establishment routes
	// we need to create authorization middleware for the establishment routes
}
