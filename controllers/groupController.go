package controllers

import (
	"net/http"
	"study-bean/database"
	"study-bean/responses"

	"github.com/gin-gonic/gin"
)



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

	// Check if that user already exists
	_, err := database.FindUserInGroup(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.InternalServerError,
		})
		return
	}

	// If yes return already in group.
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.UserAlreadyInGroup,
		})
	}

	// If no then add the user id to group
	err = database.AddUserToGroup(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.InternalServerError,
		})
		return
	}

	// return success true response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.UserAddedToGroup,
	})
}
