package database

import (
	"context"
	"fmt"
	"study-bean/initializers"
	"study-bean/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddGroupTodo(groupTodo models.GTodo, guid string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var gTodoObject models.GroupTodo
	filter := bson.M{"group_ref_id": guid}

	err := initializers.GroupTodoCollection.FindOne(ctx, filter).Decode(&gTodoObject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("1")
			// No existing entry, create a new one
			gTodoObject.Todos = []models.GTodo{groupTodo}
			gTodoObject.GroupRefID = guid
			gTodoObject.TodoCount = 1
			gTodoObject.Completed = 0

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			_, err = initializers.GroupTodoCollection.InsertOne(ctx, gTodoObject)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		fmt.Println("2")

		_, err = initializers.GroupTodoCollection.UpdateOne(ctx, filter, bson.M{
			"$addToSet": bson.M{
				"todos": groupTodo,
			},
			"$inc": bson.M{
				"todo_count": +1,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateGroupTodo(todo_id primitive.ObjectID, priority models.Priority, todo, guid string) (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"group_ref_id": guid, "todos.todo._id": todo_id}
	update := bson.M{
		"$set": bson.M{
			"todos.$.todo.todo_body": todo,
			"todos.$.todo.priority":  priority,
		},
	}
	fmt.Println(filter)

	result, err := initializers.GroupTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func ToggleGroupTodo(todo_id primitive.ObjectID, isCompleted bool, guid string) (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"group_ref_id": guid, "todos.todo._id": todo_id}
	update := bson.M{
		"$set": bson.M{
			"todos.$.todo.isCompleted": isCompleted,
		},
	}
	result, err := initializers.GroupTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return result, nil
}

func DeleteGroupTodo(todoId primitive.ObjectID, guid string) (*mongo.UpdateResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Define the filter to match the user and the todo item
	filter := bson.M{"group_ref_id": guid}

	// Define the update operation to pull the specific todo item from the todos array
	update := bson.M{
		"$pull": bson.M{
			"todos": bson.M{"todo._id": todoId},
		},
		"$inc": bson.M{
			"todo_count": -1,
		},
	}

	// Perform the update operation
	result, err := initializers.GroupTodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}