package models


type User struct {
	ID      		uint   `gorm:"primaryKey;column:id"`  
	Email           string `gorm:"column:email" json:"email"`
	Username    	string `gorm:"column:username" json:"username"`
	FullName    	string `gorm:"column:full_name" json:"full_name"`                  
	PasswordHash    string `gorm:"column:password_hash" json:"password"`
	Age 			uint   `gorm:"column:age" json:"age"`
	JwtToken        string `gorm:"column:jwt_token"`
}

type Post struct {
	ID		 uint   `gorm:"primaryKey;column:id"` 
	UserID   uint   `gorm:"column:user_id"`  
	Content  string `gorm:"column:content"`  
	ImageURL string `gorm:"column:image_url"` 
	Comments []Comment `gorm:"foreignKey:PostID"` 
}

type Comment struct {
	AuthorID uint   `gorm:"column:author_id"`  
	Content  string `gorm:"column:content"`  
	PostID   uint   `gorm:"column:post_id"`
	Post     Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserActivityLog struct {
	UserID 		uint   `gorm:"column:user_id"`  
	Description string `gorm:"column:description"`  
}
