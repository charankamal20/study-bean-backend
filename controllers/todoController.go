package controllers

import (
	"context"
	"fmt"
	"net/http"
	"study-bean/database"
	"study-bean/initializers"
	"study-bean/models"
	"study-bean/responses"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddTodo(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
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

	// Parse current time
	current_time := time.Now().UTC()

	// Create a new Todo item
	newTodo := models.Todo{
		ID:            primitive.NewObjectID(),
		Todo:          body.Todo,
		Priority:      body.Priority,
		IsCompleted:   false,
		DateCreated:   current_time,
		TimeCompleted: time.Time{},
	}

	var userTodo models.UserTodo
	filter := bson.M{"user_ref_id": user_id.(string)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = initializers.UserTodoCollection.FindOne(ctx, filter).Decode(&userTodo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No existing entry, create a new one
			userTodo = models.UserTodo{
				UserRefID: user_id.(string),
				Completed: 0,
				Count:     1,
				Todos:     []models.Todo{newTodo},
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_, err = initializers.UserTodoCollection.InsertOne(ctx, userTodo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": responses.ErrorCreateTodo,
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": responses.DatabaseError,
			})
			return
		}
	} else {
		// Existing entry found, update it
		userTodo.Todos = append(userTodo.Todos, newTodo)
		userTodo.Count++

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err = initializers.UserTodoCollection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				"todos": userTodo.Todos,
				"count": userTodo.Count,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": responses.ErrorUpdateTodo,
			})
			return
		}
	}

	// Respond
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": responses.SuccessAddTodo,
	})
}

func GetAllTodos(c *gin.Context) {

	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	todos, err := database.GetUserTodos(user_id.(string))
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"success": false,
			"message": responses.ErrorNoTodosFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    todos,
	})
}

func ToggleTodoState(context *gin.Context) {
	user_id, exists := context.Get("user_id")
	if !exists {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	todoIDFromParam := context.Param("todo_id")
	todoID, err := primitive.ObjectIDFromHex(todoIDFromParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}
	var toggledTodo struct {
		IsCompleted bool `json:"isCompleted"`
	}
	if err := context.ShouldBindJSON(&toggledTodo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}

	result, err := database.ToggleTodoState(todoID, toggledTodo.IsCompleted, user_id.(string))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	if result.MatchedCount == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": responses.ErrorTodoNotFound,
		})
		return
	}
	fmt.Println(result)
	// Respond
	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func UpdateTodo(context *gin.Context) {
	user_id, exists := context.Get("user_id")
	if !exists {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	todoIDParam := context.Param("todo_id")
	todoID, err := primitive.ObjectIDFromHex(todoIDParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}

	var updateData struct {
		Todo     string          `json:"todo"`
		Priority models.Priority `json:"priority"`
	}

	if err := context.ShouldBindJSON(&updateData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.TryAgain,
		})
		return
	}

	result, err := database.UpdateTodo(todoID, updateData.Priority, updateData.Todo, user_id.(string))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	if result.MatchedCount == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": responses.ErrorTodoNotFound,
		})
		return
	}

	fmt.Println(result)

	// Respond
	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func DeleteTodo(context *gin.Context) {
	user_id, exists := context.Get("user_id")
	if !exists {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	todoIDParam := context.Param("todo_id")
	todoID, err := primitive.ObjectIDFromHex(todoIDParam)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": responses.InvalidTodoID,
		})
		return
	}

	_, err = database.DeleteTodo(todoID, user_id.(string))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": responses.ErrorUpdateTodo,
		})
		return
	}

	// Respond
	context.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
