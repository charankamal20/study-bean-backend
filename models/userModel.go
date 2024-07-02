package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// album represents data about a record album.
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	User_ID      string             `json:"user_id" bson:"user_id,omitempty"`
	Email        string             `json:"email" bson:"email,omitempty"`
	Password     string             `json:"password" bson:"password,omitempty"`
	Username     string             `json:"username" bson:"username,omitempty"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type Priority string

const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

type Todo struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Todo          string             `json:"todo" bson:"todo,omitempty"`
	IsCompleted   bool               `json:"isCompleted" bson:"isCompleted,omitempty"`
	DateCreated   time.Time          `json:"dateCreated" bson:"dateCreated,omitempty"`
	Priority      Priority           `json:"priority" bson:"priority,omitempty"`
	TimeCompleted time.Time          `json:"timeCompleted" bson:"timeCompleted,omitempty"`
}

type UserTodo struct {
	UserRefID string `json:"user_ref_id" bson:"user_ref_id,omitempty"`
	Completed int    `json:"completed" bson:"completed,omitempty"`
	Count     int    `json:"count" bson:"count,omitempty"`
	Todos     []Todo `json:"todos" bson:"todos,omitempty"`
}
