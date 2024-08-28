package routes

import (
	"finalp2/controllers"
	"finalp2/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	_ "graded-challenge-3-celotip/docs"
    echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// Insert db to echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			c.Set("cfg", cfg)
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

	e.POST("/posts", controllers.CreatePost, jwtMiddleware)
	e.GET("/posts", controllers.GetAllPosts, jwtMiddleware)
	e.GET("/posts/:id", controllers.GetPostByID, jwtMiddleware)
	e.DELETE("/posts/:id", controllers.DeletePostByID, jwtMiddleware)
	e.POST("/comments", controllers.CreateComment, jwtMiddleware)
	e.GET("/comments/:id", controllers.GetCommentByID, jwtMiddleware)
	e.DELETE("/comments/:id", controllers.DeleteCommentByID, jwtMiddleware)
	e.GET("/activities", controllers.GetUserActivities, jwtMiddleware)
}
