package database

import (
	"context"
	"fmt"
	"study-bean/initializers"
	"study-bean/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
