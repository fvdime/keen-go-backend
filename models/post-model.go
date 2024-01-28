package models

import (
	"time"

	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// so struct is basically converter between db and go. db understand json, go does not understand json
type Post struct {
	Title      *string            `json:"title" validate:"required,min=2,max=50"`
	Body       *string            `json:"body" validate:"required,min=2"`
	Image      string             `json:"image" validate:"required,min=6"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
	User_Id    string             `json:"user_id"`
	Post_Id    string             `json:"post_id"`
}
