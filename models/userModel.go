package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// album represents data about a record album.
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	User_ID      string             `json:"user_id" bson:"user_id"`
	About        string             `json:"about" bson:"about"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Username     string             `json:"username" bson:"username"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	Groups       []string           `json:"groups" bson:"groups"`
}

type Priority string

const (
	Low    Priority = "Low"
	Medium Priority = "Medium"
	High   Priority = "High"
)

type Todo struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	TodoBody      string             `json:"todo_body" bson:"todo_body"`
	IsCompleted   bool               `json:"isCompleted" bson:"isCompleted"`
	DateCreated   time.Time          `json:"dateCreated" bson:"dateCreated"`
	Priority      Priority           `json:"priority" bson:"priority"`
	TimeCompleted time.Time          `json:"timeCompleted" bson:"timeCompleted"`
}

type UserTodo struct {
	UserRefID string `json:"user_ref_id" bson:"user_ref_id"`
	Completed int    `json:"completed" bson:"completed"`
	Count     int    `json:"count" bson:"count"`
	Todos     []Todo `json:"todos" bson:"todos"`
}
