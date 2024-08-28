package controllers

import (
	"encoding/json"
	"fmt"
	"finalp2/config"
	"finalp2/models"
	"finalp2/utils"
	"net/http"
	"regexp"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetJoke(c echo.Context, cfg *config.Config) (string, error) {
	url := "https://api.api-ninjas.com/v1/jokes"
	headers := map[string]string{
		"X-Api-Key": cfg.XAPIKey,
	}
	resp, err := utils.RequestGET(url, headers)
	if err != nil {
		return "", err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return "", err
	}
	joke := result[0]["joke"].(string)

	return joke, nil
}

type PostRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}

// CreatePost creates a new post
// @Summary Create a new post
// @Description Create a new post for the logged-in user. If content is empty, a random joke will be used.
// @Tags Posts
// @Accept json
// @Produce json
// @Param postRequest body PostRequest true "Post Request"
// @Success 201 {object} map[string]interface{} "Post created successfully"
// @Failure 400 {object} utils.APIError "Invalid request"
// @Failure 500 {object} utils.APIError "Failed to create post"
// @Security ApiKeyAuth
// @Router /posts [post]
func CreatePost(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	cfg := c.Get("cfg").(*config.Config)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var postRequest PostRequest
	if err := c.Bind(&postRequest); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid request"))
	}

	// Validate Image URL
	regex := `^(http|https)://[^\s/$.?#].[^\s]*$`
	validURL, err := regexp.MatchString(regex, postRequest.ImageURL)
	if err != nil || !validURL {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid image URL"))
	}

	// If content is empty, get a random joke
	if postRequest.Content == "" {
		joke, err := GetJoke(c, cfg)
		if err != nil {
			return utils.HandleError(c, utils.NewInternalError("Failed to get joke"))
		}
		postRequest.Content = joke
	}

	// Create Post
	post := models.Post{
		UserID:   uint(userID),
		Content:  postRequest.Content,
		ImageURL: postRequest.ImageURL,
	}
	if err := db.Create(&post).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to create post"))
	}

	logEntry := models.UserActivityLog{
		UserID: uint(userID),
		Description: fmt.Sprintf("user create new POST with ID %d", post.ID),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to log user activity"))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post created successfully",
		"post":    post,
	})
}

// GetAllPosts returns all posts
// @Summary Get all posts
// @Description Get all posts stored in the database
// @Tags Posts
// @Produce json
// @Success 200 {array} models.Post "List of all posts"
// @Failure 500 {object} utils.APIError "Error fetching posts"
// @Security ApiKeyAuth
// @Router /posts [get]
func GetAllPosts(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)

	var posts []models.Post
	if err := db.Find(&posts).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching posts"))
	}

	return c.JSON(http.StatusOK, posts)
}

// GetPostByID returns a post by ID
// @Summary Get a post by ID
// @Description Get a specific post by its ID, including its comments
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post "Post details"
// @Failure 404 {object} utils.APIError "Post not found"
// @Security ApiKeyAuth
// @Router /posts/{id} [get]
func GetPostByID(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	param := c.Param("id")
	postID, _ := strconv.Atoi(param)

	var post models.Post
	if err := db.Preload("Comments").First(&post, postID).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Post not found"))
	}

	return c.JSON(http.StatusOK, post)
}

// DeletePostByID deletes a post by ID
// @Summary Delete a post by ID
// @Description Delete a specific post by its ID. Only the owner of the post can delete it.
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post deleted successfully"
// @Failure 404 {object} utils.APIError "Post not found"
// @Failure 401 {object} utils.APIError "You are not authorized to delete this post"
// @Security ApiKeyAuth
// @Router /posts/{id} [delete]
func DeletePostByID(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)
	param := c.Param("id")
	postID, _ := strconv.Atoi(param)

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Post not found"))
	}

	// Check if the user is the owner of the post
	if post.UserID != uint(userID) {
		return utils.HandleError(c, utils.NewUnauthorizedError("You are not authorized to delete this post"))
	}

	if err := db.Delete(&post).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error deleting post"))
	}

	logEntry := models.UserActivityLog{
		UserID: uint(userID),
		Description: fmt.Sprintf("user delete POST with ID %d", post.ID),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to log user activity"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Post deleted successfully",
		"post":    post,
	})
}

type CommentRequest struct {
	Content string `json:"content"`
	PostID  uint    `json:"post_id"`
}

// CreateComment creates a new comment on a post
// @Summary Create a new comment
// @Description Create a new comment on a specific post
// @Tags Comments
// @Accept json
// @Produce json
// @Param commentRequest body CommentRequest true "Comment Request"
// @Success 201 {object} map[string]interface{} "Comment created successfully"
// @Failure 400 {object} utils.APIError "Invalid request"
// @Failure 500 {object} utils.APIError "Failed to create comment"
// @Security ApiKeyAuth
// @Router /comments [post]
func CreateComment(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var commentRequest CommentRequest
	if err := c.Bind(&commentRequest); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid request"))
	}

	comment := models.Comment{
		AuthorID: uint(userID),
		PostID:   commentRequest.PostID,
		Content:  commentRequest.Content,
	}
	if err := db.Create(&comment).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to create comment"))
	}

	logEntry := models.UserActivityLog{
		UserID: uint(userID),
		Description: fmt.Sprintf("user create new comment in Post ID %d", comment.PostID),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to log user activity"))
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetCommentByID returns a comment by ID
// @Summary Get a comment by ID
// @Description Get a specific comment by its ID, including the associated post and author details
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.Comment "Comment details"
// @Failure 404 {object} utils.APIError "Comment not found"
// @Security ApiKeyAuth
// @Router /comments/{id} [get]
func GetCommentByID(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	param := c.Param("id")
	commentID, _ := strconv.Atoi(param)
	
	var comment models.Comment
	if err := db.Preload("Post").Preload("Author").First(&comment, commentID).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Comment not found"))
	}

	return c.JSON(http.StatusOK, comment)
}

// DeleteCommentByID deletes a comment by ID
// @Summary Delete a comment by ID
// @Description Delete a specific comment by its ID. Only the author of the comment can delete it.
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{} "Comment deleted successfully"
// @Failure 404 {object} utils.APIError "Comment not found"
// @Failure 401 {object} utils.APIError "You are not authorized to delete this comment"
// @Security ApiKeyAuth
// @Router /comments/{id} [delete]
func DeleteCommentByID(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)
	param := c.Param("id")
	commentID, _ := strconv.Atoi(param)

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Comment not found"))
	}

	// Check if the user is the owner of the comment
	if comment.AuthorID != uint(userID) {
		return utils.HandleError(c, utils.NewUnauthorizedError("You are not authorized to delete this comment"))
	}

	if err := db.Delete(&comment).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error deleting comment"))
	}

	logEntry := models.UserActivityLog{
		UserID: uint(userID),
		Description: fmt.Sprintf("user delete comment in POST ID %d", comment.PostID),
	}
	if err := db.Create(&logEntry).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to log user activity"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Comment deleted successfully",
		"comment": comment,
	})
}

// GetUserActivities returns the activity logs of the logged-in user
// @Summary Get user activities
// @Description Get the activity logs of the logged-in user
// @Tags Activities
// @Produce json
// @Success 200 {array} models.UserActivityLog "List of user activities"
// @Failure 500 {object} utils.APIError "Error fetching user activities"
// @Security ApiKeyAuth
// @Router /activities [get]
func GetUserActivities(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var activities []models.UserActivityLog
	if err := db.Where("user_id = ?", userID).Find(&activities).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching user activities"))
	}

	return c.JSON(http.StatusOK, activities)
}




