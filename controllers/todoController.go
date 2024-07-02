package controllers

import (
	"net/http"
	"study-bean/database"
	"study-bean/models"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func AddTodo(context *gin.Context) {

	user_id, exists := context.Get("user_id")

	if !exists {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	todos := []models.Todo{
		{
			ID:            primitive.NewObjectID(),
			Todo:          "Hello",
			Priority:      models.High,
			IsCompleted:   false,
			DateCreated:   time.Unix(time.Now().Unix(), 0),
			TimeCompleted: time.Unix(time.Now().Unix(), 0),
		},
	}

	userTodo := models.UserTodo {
		Completed: 0,
		Count: 1,
		Todos: todos,
		UserRefID: user_id.(string),
	}

	err := database.AddTodoToUser(userTodo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": "Please Try Again",
		})
		return
	}

	// Respond
	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"todo": userTodo,
	})
}