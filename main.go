package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type user struct {
    ID     string  `json:"id"`
    Username  string  `json:"username"`
    Password string  `json:"password"`
}
// albums slice to seed record album data.
var users = []user{
    {ID: "1", Username: "Blue Train", Password: "John Coltrane"},
    {ID: "2", Username: "Jeru", Password: "Gerry Mulligan"},
    {ID: "3", Username: "Sarah Vaughan and Clifford Brown", Password: "Sarah Vaughan"},
}

func addUser(c *gin.Context) {
	var newUser user

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

// getAlbums responds with the list of all albums as JSON.
func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func getUserById(id string) (*user, error) {
    for i, u := range users {
        if u.ID == id {
            return &users[i], nil
        }
    }

    return nil, errors.New("user not found")
}

func getUser(context *gin.Context) {
    id := context.Param("id")

    user, err := getUserById(id)

    if err != nil {
        context.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
        return
    }

    context.JSON(http.StatusOK, user)
}

func updatePassword(context *gin.Context) {
    id := context.Param("id")

    user, err := getUserById(id)

    if err != nil {
        context.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
        return
    }

    var body struct {
        Password string `json:"password" binding:"required"`
    }

    if err := context.ShouldBindJSON(&body); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
    }

    user.Password = body.Password
    context.JSON(http.StatusOK, user)
}

func main() {
    router := gin.Default()
    router.GET("/user", getUsers)
	router.POST("/user", addUser)
    router.GET("/user/:id", getUser)
    router.PATCH("/user/:id", updatePassword)
    router.Run("localhost:8080")
}