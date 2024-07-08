package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Group struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	GroupName        string             `json:"group_name" bson:"group_name"`
	GroupDescription string             `json:"group_description" bson:"group_description"`
	GroupPhoto       string             `json:"group_photo" bson:"group_photo"`
	GroupID          string             `json:"group_id" bson:"group_id"`
	AdminID          string             `json:"admin_id" bson:"admin_id"`
	NumberOfMembers  int                `json:"number_of_members" bson:"number_of_members"`
	Members          []string           `json:"members" bson:"members"`
	Banned           []string           `json:"banned" bson:"banned"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type GTodo struct {
	Todo
	Creator string `json:"creator" bson:"creator"`
}

type GroupTodo struct {
	GroupRefID string  `json:"group_ref_id" bson:"group_ref_id"`
	Completed  int     `json:"completed" bson:"completed"`
	TodoCount  int     `json:"todo_count" bson:"todo_count"`
	Todos      []GTodo `json:"todos" bson:"todos"`
}
