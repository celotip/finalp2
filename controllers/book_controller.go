package controllers

import (
	"finalp2/helper"
	"finalp2/models"
	"finalp2/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TopupRequest struct {
	Amount  uint `json:"amount"`
}

// Topup adds a deposit to the user's account
// @Summary Add deposit to user account
// @Description Add a specified amount to the user's deposit
// @Tags Users
// @Accept json
// @Produce json
// @Param topupRequest body TopupRequest true "Topup Request"
// @Success 200 {object} map[string]interface{} "Deposit added successfully"
// @Failure 400 {object} utils.APIError "Invalid request"
// @Failure 500 {object} utils.APIError "Failed to update user"
// @Security ApiKeyAuth
// @Router /topup [post]
func Topup(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var topupRequest TopupRequest
	if err := c.Bind(&topupRequest); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid request"))
	}

	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching user"))
	}

	user.Deposit += topupRequest.Amount

	// Save method updates the entire product record in the database
	if err := db.Save(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to update user"))
	}

	outputMessage := fmt.Sprintf("Deposit added. Current deposit amount: %d", user.Deposit)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": outputMessage,
	})
}

type BookResponse struct {
    ID         uint       `json:"id"`
    Title      string     `json:"title"`
    Author     string `json:"author"`
    Category   string `json:"category"`
}

// GetAllBooks returns all books
// @Summary Get all books
// @Description Get all books stored in the database
// @Tags Books
// @Produce json
// @Success 200 {array} models.Book "List of all books"
// @Failure 500 {object} utils.APIError "Error fetching books"
// @Security ApiKeyAuth
// @Router /books [get]
func GetAllBooks(c echo.Context) error {
    db := c.Get("db").(*gorm.DB)

    var books []models.Book
    if err := db.Preload("Author").Preload("Category").Find(&books).Error; err != nil {
        return utils.HandleError(c, utils.NewInternalError("Error fetching books"))
    }

    var bookResponses []BookResponse
    for _, book := range books {
        bookResponses = append(bookResponses, BookResponse{
            ID:       book.ID,
            Title:    book.Title,
            Author:   book.Author.FirstName + " " + book.Author.LastName,
            Category: book.Category.Name,
        })
    }

    return c.JSON(http.StatusOK, bookResponses)
}

// GetBookById returns a specific book by ID
// @Summary Get book by ID
// @Description Get details of a book by its ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book "Book details"
// @Failure 404 {object} utils.APIError "Book not found"
// @Failure 500 {object} utils.APIError "Error fetching book"
// @Security ApiKeyAuth
// @Router /books/{id} [get]
func GetBookById(c echo.Context) error {
    db := c.Get("db").(*gorm.DB)
    param := c.Param("id")
    bookID, _ := strconv.Atoi(param)

    var book models.Book
    if err := db.Preload("Author").Preload("Category").Where("book_id = ?", bookID).First(&book).Error; err != nil {
        return utils.HandleError(c, utils.NewNotFoundError("Book not found"))
    }

    response := BookResponse{
        ID:       book.ID,
        Title:    book.Title,
        Author:   book.Author.FirstName + " " + book.Author.LastName,
        Category: book.Category.Name,
    }

    return c.JSON(http.StatusOK, response)
}

// GetCart returns the user's cart
// @Summary Get user's cart
// @Description Get the current items in the user's cart
// @Tags Cart
// @Produce json
// @Success 200 {array} models.Book "List of books in the cart"
// @Failure 500 {object} utils.APIError "Error fetching cart"
// @Security ApiKeyAuth
// @Router /cart [get]
func GetCart(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var out []models.Book
	var carts []models.Cart
	if err := db.Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching cart"))
	}

	for _, cart := range carts {
		var product models.Book
		if err := db.Where("book_id = ?", cart.BookID).First(&product).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Error fetching cart"))
		}
		out= append(out, product)
	}

	return c.JSON(http.StatusOK, out)
}

type CartInput struct {
	BookID      		uint   `json:"book_id"`  
}

// AddCart adds a book to the user's cart
// @Summary Add book to cart
// @Description Add a book to the user's cart by its ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param cartInput body CartInput true "Cart Input"
// @Success 200 {object} CartInput "Book added to cart"
// @Failure 400 {object} utils.APIError "Invalid input"
// @Failure 500 {object} utils.APIError "Failed to add book to cart"
// @Security ApiKeyAuth
// @Router /cart [post]
func AddCart(c echo.Context) error {
	cart := new(models.Cart)
	cartInp := new(CartInput)

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	cart.UserID = uint(userID)

	if err := c.Bind(cartInp); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	cart.BookID = cartInp.BookID

	db := c.Get("db").(*gorm.DB)

	if err := db.Create(&cart).Error; err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Failed to add product to cart."))
	}

	return c.JSON(http.StatusOK, cartInp)
}

// DeleteCart removes a book from the user's cart
// @Summary Remove book from cart
// @Description Remove a book from the user's cart by its ID
// @Tags Cart
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{} "Cart deleted successfully"
// @Failure 404 {object} utils.APIError "Cart not found"
// @Failure 500 {object} utils.APIError "Error deleting cart"
// @Security ApiKeyAuth
// @Router /cart/{id} [delete]
func DeleteCart(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var cart models.Cart
	if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Post not found"))
	}

	if err := db.Delete(&cart).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error deleting post"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Cart deleted successfully",
	})
}

type OutOrder struct {
	OrderID      		uint   `json:"rental_id"`  
	TotalPrice			uint   `json:"total_price"` 
	Date    			*time.Time `json:"date"`
	Status				string `json:"status"`
	Books				[]models.Book `json:"books"`  
}

// GetRent returns all rentals made by the user
// @Summary Get user rentals
// @Description Get all rentals made by the user
// @Tags Rentals
// @Produce json
// @Success 200 {array} OutOrder "List of user rentals"
// @Failure 500 {object} utils.APIError "Error fetching rentals"
// @Security ApiKeyAuth
// @Router /rentals [get]
func GetRent(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var out []OutOrder
	var books []models.Book
	var orders []models.Rental
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching orders"))
	}

	for _, order := range orders {
		var orderItems []models.RentalDetail
		if err := db.Where("rental_id = ?", order.ID).Find(&orderItems).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Error fetching rent detail"))
		}
		for _, orderItem := range orderItems {
			var product models.Book
			if err := db.Where("book_id = ?", orderItem.BookID).First(&product).Error; err != nil {
				return utils.HandleError(c, utils.NewInternalError("Error fetching book"))
			}
			books = append(books, product)
		}
		var outOrder OutOrder
		outOrder.OrderID = order.ID
		outOrder.TotalPrice = order.TotalPrice
		outOrder.Date = order.RentalDate
		outOrder.Status = order.RentalStatus
		outOrder.Books = books
		out = append(out, outOrder)
	}

	return c.JSON(http.StatusOK, out)
}

// AddOrder creates a new order from the user's cart
// @Summary Create a new order
// @Description Create a new order from the items in the user's cart. The cart will be cleared after the order is created.
// @Tags Orders
// @Produce json
// @Success 200 {object} map[string]interface{} "Order created successfully"
// @Failure 400 {object} utils.APIError "Cart is empty, cannot create order"
// @Failure 500 {object} utils.APIError "Failed to create order or order items"
// @Security ApiKeyAuth
// @Router /orders [post]
func AddOrder(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	var carts []models.Cart
	if err := db.Where("user_id = ?", uint(userID)).Find(&carts).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching cart"))
	}

	if len(carts) == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Cart is empty, cannot create order"))
	}

	// Create a new order
	order := models.Rental{
		UserID:     uint(userID),
		TotalPrice: 0, // Will be calculated below
		RentalDate:       nil,
		RentalStatus: "created",
	}

	// Calculate total price and prepare order items
	var totalPrice uint
	var orderItems []models.RentalDetail

	for _, cartItem := range carts {
		var product models.Book
		if err := db.Where("book_id = ?", cartItem.BookID).First(&product).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Error fetching book for cart item"))
		}

		orderItem := models.RentalDetail{
			RentalID: order.ID,
			BookID:  cartItem.BookID,
			Returned:     false,
		}
		orderItems = append(orderItems, orderItem)

		// Calculate the total price
		totalPrice += product.Price
	}

	// Set the total price of the order
	order.TotalPrice = totalPrice

	// Save the order
	if err := db.Create(&order).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to create order"))
	}

	// Set the OrderID in each order item and save them
	for i := range orderItems {
		orderItems[i].RentalID = order.ID
		if err := db.Create(&orderItems[i]).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Failed to create order item"))
		}
	}

	// Clear the cart after creating the order
	if err := db.Where("user_id = ?", uint(userID)).Delete(&models.Cart{}).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to clear cart after creating order"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order created successfully",
		"order_id": order.ID,
		"total_price": order.TotalPrice,
	})
}

func Pay(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(float64)

	// Extract order ID from the request body or URL parameters
	orderID, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid order ID"))
	}

	// Fetch the order from the database
	var order models.Rental
	if err := db.Where("rental_id = ?", orderID).First(&order).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("Order not found"))
	}

	if order.UserID != uint(userID) {
		return utils.HandleError(c, utils.NewUnauthorizedError("You are not authorized to pay for this order"))
	}

	// Check if the order has already been paid
	if order.RentalStatus == "paid" {
		return utils.HandleError(c, utils.NewBadRequestError("Order is already paid"))
	}

	// Proceed with the payment process (this is where you integrate with a payment gateway or handle payment logic)
	// Assuming payment is successful

	// Update the order status to 'paid'
	order.RentalStatus = "paid"
	if err := db.Save(&order).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Failed to update order status to paid"))
	}

	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return utils.HandleError(c, utils.NewNotFoundError("User not found"))
	}

	var orderItems []models.RentalDetail
	if err := db.Where("rental_id = ?", orderID).First(&orderItems).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error fetching ordered books"))
	}

	var books []models.Book
	for _, orderItem := range orderItems {
		var product models.Book
		if err := db.Where("book_id = ?", orderItem.BookID).First(&product).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Error fetching book for cart item"))
		}
		books = append(books, product)
	}


	// Optionally, generate an invoice or payment confirmation
	invoiceRes, err := helper.CreateInvoice(order, user, books)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error while creating invoice")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment successful",
		"invoice": invoiceRes,
		"order_id": order.ID,
		"status":  order.RentalStatus,
	})
}