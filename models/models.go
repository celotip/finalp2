package models

import "time"


type User struct {
	ID      		uint   `gorm:"primaryKey;column:user_id"`  
	Email           string `gorm:"column:email" json:"email"`
	Password    	string `gorm:"column:password_hash" json:"password"`
	FirstName    	string `gorm:"column:first_name" json:"first_name"`
	LastName    	string `gorm:"column:last_name" json:"last_name"`
	Birth_date		string `gorm:"column:birth_date" json:"birth_date"`
	Address			string `gorm:"column:address" json:"address"`
	Contact			string `gorm:"column:contact_no" json:"contact_no"`
	Deposit			uint   `gorm:"column:deposit" json:"deposit"`
	JwtToken        string `gorm:"column:jwt_token"`

}

type Book struct {
	ID      		uint   `gorm:"primaryKey;column:book_id"`
	Title		  	string `gorm:"column:title" json:"title"`  
	AuthorID    	uint `gorm:"column:author_id" json:"author"`
	CategoryID		uint `gorm:"column:category_id" json:"category"`
	Author     		Author
    Category   		Category
	Isbn			string `gorm:"column:isbn" json:"ISBN"`
	Stock			uint   `gorm:"column:stock" json:"stock"`
	Price			uint   `gorm:"column:price" json:"price"`
	ReadingDays		uint   `gorm:"column:reading_days" json:"reading_days"`
}

type Author struct {
	ID         uint   `gorm:"primaryKey;column:author_id"`
	FirstName  string `gorm:"column:first_name" json:"first_name"`
	LastName   string `gorm:"column:last_name" json:"last_name"`
	Nationality string `gorm:"column:nationality" json:"nationality"`
	BirthDate  string `gorm:"column:birth_date" json:"birth_date"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey;column:category_id"`
	Name string `gorm:"column:name" json:"name"`
}

type Rental struct {
	ID          uint   `gorm:"primaryKey;column:rental_id"`
	UserID      uint   `gorm:"column:user_id" json:"user_id"`
	RentalDate  *time.Time `gorm:"column:rental_date" json:"rental_date"`
	RentalStatus string `gorm:"column:rental_status" json:"rental_status"`
	TotalPrice   uint `gorm:"column:total_price" json:"total_price"`
}

type RentalDetail struct {
	ID         uint   `gorm:"primaryKey;column:rental_detail_id"`
	RentalID   uint   `gorm:"column:rental_id" json:"rental_id"`
	BookID     uint   `gorm:"column:book_id" json:"book_id"`
	Returned   bool   `gorm:"column:returned" json:"returned"`
}

type Payment struct {
	ID            uint   `gorm:"primaryKey;column:payment_id"`
	RentalID      uint   `gorm:"column:rental_id" json:"rental_id"`
	PaymentDate   string `gorm:"column:payment_date" json:"payment_date"`
	PaymentAmount float64 `gorm:"column:payment_amount" json:"payment_amount"`
}

type Cart struct {
	UserID				uint   `gorm:"column:user_id" json:"user_id"`
	BookID      		uint   `gorm:"column:book_id" json:"book_id"`  
}
