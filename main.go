package main

import (
	"finalp2/config"
	"finalp2/middlewares"
	"finalp2/routes"
	"finalp2/utils"
	"net/http"

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

// @title Social Media
// @version 1.0
// @description Ini adalah API social media untuk membuat post, delete post, comment, dll
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

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()
	
	cfg := config.LoadConfig()

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

	routes.SetupRoutes(e, db, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	e.Logger.Fatal(e.Start(":"+port))

	// // Mengambil objek database SQL dari GORM
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("failed to get database: %v", err)
	// }

	// // Mengatur pooling koneksi
	// sqlDB.SetMaxIdleConns(10)                 // Jumlah maksimal koneksi idle
	// sqlDB.SetMaxOpenConns(100)                // Jumlah maksimal koneksi yang terbuka sekaligus
	// sqlDB.SetConnMaxLifetime(time.Hour)       // Durasi maksimal sebuah koneksi (dalam jam)
	// sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Durasi maksimal sebuah koneksi idle (dalam menit)

	// // Menutup koneksi database saat aplikasi berhenti
	// defer func() {
	// 	err := sqlDB.Close()
	// 	if err != nil {
	// 		log.Fatalf("failed to close database: %v", err)
	// 	}
	// }()
}