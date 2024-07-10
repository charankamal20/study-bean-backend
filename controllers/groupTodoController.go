package controllers

import (
	"fmt"
	"net/http"
	"study-bean/database"
	"study-bean/models"
	"study-bean/responses"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddGroupTodo(c *gin.Context) {

	// Get User id and email from middleware
	user_id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}

	guid := c.Param("guid")
	fmt.Println("provider", guid)
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.GroupIDMissing,
		})
		return
	}

	var body struct {
		Todo     string
		Priority models.Priority
	}

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TodoMissing,
		})
		return
	}

	var groupTodo models.GTodo

	groupTodo.Creator = user_id.(string)
	groupTodo.DateCreated = time.Now()
	groupTodo.ID = primitive.NewObjectID()
	groupTodo.IsCompleted = false
	groupTodo.Priority = body.Priority
	groupTodo.TodoBody = body.Todo

	err = database.AddGroupTodo(groupTodo, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorCreateTodo,
			"error":   err,
		})
		return
	}

	// Respond
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.SuccessAddTodo,
		"todo":    groupTodo,
	})
}

func UpdateGroupTodo(c *gin.Context) {
	todoIDParam := c.Query("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDParam)
	if todoIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.GroupIDMissing,
		})
		return
	}

	var updateData struct {
		Todo     string          `json:"todo"`
		Priority models.Priority `json:"priority"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TodoMissing,
		})
		return
	}
	if updateData.Todo == "" && updateData.Priority == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TodoMissing,
		})
		return
	}

	result, err := database.UpdateGroupTodo(todoID, updateData.Priority, updateData.Todo, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": responses.ErrorTodoNotFound,
		})
		return
	}

	fmt.Println(result)

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})

}

func ToggleGroupTodo(c *gin.Context) {
	todoIDParam := c.Query("id")
	todoID, _ := primitive.ObjectIDFromHex(todoIDParam)
	if todoIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.GroupIDMissing,
		})
		return
	}

	var updateData struct {
		IsCompleted bool `json:"isCompleted"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}

	result, err := database.ToggleGroupTodo(todoID, updateData.IsCompleted, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": responses.ErrorTodoNotFound,
		})
		return
	}

	fmt.Println(result)

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func DeleteGroupTodo(c *gin.Context) {
	todoIDParam := c.Query("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDParam)
	if todoIDParam == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}

	guid := c.Param("guid")
	if guid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.GroupIDMissing,
		})
		return
	}


	_, err = database.DeleteGroupTodo(todoID, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}