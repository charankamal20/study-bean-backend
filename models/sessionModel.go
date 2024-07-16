package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	SessionName        string             `json:"session_name" bson:"session_name"`
	SessionDescription string             `json:"session_description" bson:"session_description"`
	SessionPhoto       string             `json:"session_photo" bson:"session_photo"`
	SessionID          string             `json:"session_id" bson:"session_id"`
	AdminID            string             `json:"admin_id" bson:"admin_id"`
	NumberOfMembers    int                `json:"number_of_members" bson:"number_of_members"`
	Members            []string           `json:"members" bson:"members"`
	Banned             []string           `json:"banned" bson:"banned"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
	TempMemberList     []string           `json:"temp_member_list" bson:"temp_member_list"`
	ExpiresAfter       time.Time          `json:"expiresAfter" bson:"expiresAfter"`
	SessionJoinLink    string             `json:"session_join_link" bson:"session_join_link"`
}

type STodo struct {
	Todo
	Creator string `json:"creator" bson:"creator"`
}

type SessionTodo struct {
	SessionRefID string  `json:"session_ref_id" bson:"session_ref_id"`
	Completed    int     `json:"completed" bson:"completed"`
	TodoCount    int     `json:"todo_count" bson:"todo_count"`
	Todos        []STodo `json:"todos" bson:"todos"`
}
