package routes

import (
	"finalp2/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	_ "finalp2/docs"
    echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Insert db to echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	// Without jwt tokens
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/users/register", controllers.RegisterUser)
	e.POST("/users/login", controllers.LoginUser)

	// With jwt tokens
	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		Claims:     jwt.MapClaims{},
	})

	e.GET("/books/all", controllers.GetAllBooks, jwtMiddleware)
	e.POST("/users/topup", controllers.Topup, jwtMiddleware)
	e.POST("/users/rent", controllers.AddCart, jwtMiddleware)
	e.GET("/users/carts", controllers.GetCart, jwtMiddleware)
	e.DELETE("/users/carts", controllers.DeleteCart, jwtMiddleware)
	e.POST("/users/checkout", controllers.AddOrder, jwtMiddleware)
	// e.POST("/users/pay", controllers.Pay, jwtMiddleware)
	// e.POST("/users/return", controllers.Return, jwtMiddleware)
	e.GET("/users/rent-history", controllers.GetRent, jwtMiddleware)
}
