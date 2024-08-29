package main

import (
	"context"
	"finalp2/config"
	"finalp2/middlewares"
	"finalp2/routes"
	"finalp2/utils"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func customHTTPErrorHandler(err error, c echo.Context, logger *logrus.Logger) {
	//default status code
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	logger.Info("Memasuki customHTTPErrorHandler")

	//type assertion to check if the error is an APIError
	if apiErr, ok := err.(*utils.APIError); ok {
		code = apiErr.Code
		message = apiErr.Message
	} else if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	logger.WithFields(logrus.Fields{
		"status": code,
		"error":  err.Error(),
	}).Error("HTTP Error terjadi")

	//send error response
	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})

	logger.Info("Keluar dari customHTTPErrorHandler")
}

// @title Book Rental
// @version 1.0
// @description Ini adalah API book rental untuk merental buku dari pilihan buku dan kategori yang tersedia
// @termOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@warungapi.com
// @host localhost:8080
// @basePath /

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	db := config.InitDB()

	e := echo.New()

	//Initialize Logrus Logger
	logger := logrus.New()

	logger.Info("Aplikasi dimulai")

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log Level using Logrus
	logger.SetLevel(logrus.WarnLevel)

	//Output Destination
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file) //mencatat log ke file
	} else {
		logger.Info("Failed to log to file, using default sttder")
	}

	//Custom HTTPErrorHandler
	// e.HTTPErrorHandler = customHTTPErrorHandler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		customHTTPErrorHandler(err, c, logger)
	}

	// Middleware logger
	e.Use(middleware.Logger())

	// Middleware recovery (optional, handles panics)
	e.Use(middleware.Recover())

	//middleware logrus
	e.Use(middlewares.LogrusMiddleware(logger))

	routes.SetupRoutes(e, db)

	// Handle graceful shutdown
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server...")

	// Properly shutdown the Echo server
	if err := e.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}

	// Close the database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from GORM: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close the database connection: %v", err)
	}

	log.Println("Server shut down gracefully.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	e.Logger.Fatal(e.Start(":"+port))
}