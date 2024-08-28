package controllers

import (
	"graded-challenge-3-celotip/models"
	"graded-challenge-3-celotip/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInput struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}

type UserOutput struct {
	ID      		uint   `json:"user_id"`  
	Email           string `json:"email"`
	Username    	string `json:"username"`
	FullName    	string `json:"full_name"`                  
	Age 			uint   `json:"age"`
}


// @Summary Register a new user
// @Description Registers a new user with an email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param   userInput  body  UserInput  true  "User Registration Input"
// @Success 200 {object} models.User "User object containing ID, Email, and other fields"
// @Failure 400 {object} utils.APIError "Invalid input or failed to create user"
// @Router /users/register [post]
func RegisterUser(c echo.Context) error {
	// Bind the input JSON to User
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 16)
	user.PasswordHash = string(hashedPassword)

	db := c.Get("db").(*gorm.DB)

	// Save the user to db
	if err := db.Create(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Failed to create user. Username or email already exists"))
	}

	out := new(UserOutput)
	out.ID = user.ID
	out.Username = user.Username
	out.Email = user.Email
	out.FullName = user.FullName
	out.Age = user.Age

	return c.JSON(http.StatusOK, out)
}


// @Summary Login a user
// @Description Logs in a user with an email and password, returns a JWT token
// @Tags users
// @Accept  json
// @Produce  json
// @Param   userInput  body  UserInput  true  "User Login Input"
// @Success 200 {object} map[string]string "JWT Token"
// @Failure 400 {object} utils.APIError "Invalid input or incorrect password"
// @Failure 404 {object} utils.APIError "Email not found"
// @Failure 500 {object} utils.APIError "Failed to generate token"
// @Router /users/login [post]
func LoginUser(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	input := new(UserInput)
	if err := c.Bind(input); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	dbUser := new(models.User)
	result := db.Where("username = ?", input.Username).First(&dbUser)
	if result.Error != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Email not found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(input.Password)); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Incorrect password"))
	}

	// JWT Token
	claims := jwt.MapClaims{
		"user_id":    dbUser.ID,
		"username":      dbUser.Username,
		"exp":        time.Now().Add(time.Hour * 72).Unix(), // Token expiry
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to generate token"))
	}
	dbUser.JwtToken = signedToken

	db.Save(&dbUser)
	return c.JSON(http.StatusOK, echo.Map{
		"token": dbUser.JwtToken,
	})
}
