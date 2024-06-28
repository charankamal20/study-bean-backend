package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// album represents data about a record album.
type User struct {
    ID          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    Email       string              `json:"email" bson:"email,omitempty"`
    Password    string              `json:"password" bson:"password,omitempty"`
    Username    string              `json:"username" bson:"username,omitempty"`
}


// albums slice to seed record album data.
var Users = []User{
    {Email: "Blue Train", Password: "John Coltrane", Username: "blue_train"},
    {Email: "Jeru", Password: "Gerry Mulligan", Username: "jeru"},
    {Email: "Sarah Vaughan and Clifford Brown", Password: "Sarah Vaughan", Username: "sarah_vaughan"},
}
