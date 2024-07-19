package database

import (
	"context"
	"study-bean/initializers"
	"study-bean/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserTodos(user_id string) (models.UserTodo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"user_ref_id": user_id,
	}

	var results models.UserTodo
	err := initializers.UserTodoCollection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return models.UserTodo{}, err
	}

	return results, nil
}

func UpdateTodo(todo_id primitive.ObjectID, priority models.Priority, todo, user_id string) (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"user_ref_id": user_id, "todos._id": todo_id}
	update := bson.M{
		"$set": bson.M{
			"todos.$.todo_body":     todo,
			"todos.$.priority": priority,
		},
	}

	result, err := initializers.UserTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ToggleTodoState(todo_id primitive.ObjectID, isCompleted bool, user_id string) (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"user_ref_id": user_id, "todos._id": todo_id}
	update := bson.M{
		"$set": bson.M{
			"todos.$.isCompleted": isCompleted,
		},
	}

	result, err := initializers.UserTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func DeleteTodo(todo_id primitive.ObjectID, user_id string) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Define the filter to match the user and the todo item
	filter := bson.M{"user_ref_id": user_id}

	// Define the update operation to pull the specific todo item from the todos array
	update := bson.M{
		"$pull": bson.M{
			"todos": bson.M{"_id": todo_id},
		},
		"$inc": bson.M{
			"count": -1,
		},
	}

	// Perform the update operation
	result, err := initializers.UserTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
