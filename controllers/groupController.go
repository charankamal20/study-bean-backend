package controllers

import (
	"fmt"
	"net/http"
	"study-bean/database"
	"study-bean/models"
	"study-bean/responses"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ! BUG : if user is deleted and that token is used to send a request. the request is accepted.
func CreateGroup(c *gin.Context) {

	// Get User/Admin data
	// Get User id and email from middleware
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}
	fmt.Println("USER ID", user_id.(string))

	// Get group data from req body
	var body struct {
		GroupDescription string `json:"group_description"`
		GroupName        string `json:"group_name"`
		GroupPhoto       string `json:"group_photo"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}
	fmt.Println("BODY", body)

	// Create a new Group instance
	var newGroup models.Group

	newGroup.AdminID = user_id.(string)
	newGroup.ID = primitive.NewObjectID()
	newGroup.GroupDescription = body.GroupDescription
	newGroup.GroupName = body.GroupName
	newGroup.GroupPhoto = body.GroupPhoto
	newGroup.UpdatedAt = time.Now()
	newGroup.NumberOfMembers = 1
	newGroup.Members = []string{user_id.(string)}
	newGroup.Banned = []string{}
	newGroup.GroupID = newGroup.ID.Hex()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.AddGroupInUser(newGroup.GroupID, user_id.(string))
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := database.CreateGroup(&newGroup)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": responses.InternalServerError,
			})
			return
		}
	}

	// return success true response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.NewGroupCreated,
	})
}

func AddUserToGroup(c *gin.Context) {

	// Get User id and email from middleware
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}
	fmt.Println("USER ID", user_id.(string))

	// Get Group ID from query params
	guid := c.Query("guid")
	fmt.Println("GUID", guid)

	// Check if user already in group
	user, err := database.CheckUserExistInGroup(user_id.(string), guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.InternalServerError,
		})
		return
	}

	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserAlreadyInGroup,
		})
		return
	}

	// Add user to group
	err = database.AddUserToGroup(user_id.(string), guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.FailedToAddUser,
		})
		return
	}

	// return success true response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.UserAddedToGroup,
	})
}
