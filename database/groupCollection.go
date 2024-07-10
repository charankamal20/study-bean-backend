package database

import (
	"context"
	"study-bean/initializers"
	"study-bean/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUserToGroup(user_id string, guid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Define filter and update for the group
	filterGroup := bson.M{"group_id": guid}
	updateGroup := bson.M{
		"$addToSet": bson.M{"members": user_id},
		"$inc":      bson.M{"number_of_members": 1},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	// Define filter and update for the user
	filterUser := bson.M{"user_id": user_id}
	updateUser := bson.M{
		"$addToSet": bson.M{"groups": guid},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	// Update group
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := initializers.GroupCollection.UpdateOne(ctx, filterGroup, updateGroup)
		if err != nil {
			errChan <- err
		}
	}()

	// Update user
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := initializers.UserCollection.UpdateOne(ctx, filterUser, updateUser)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func FindUserInGroup(user_id string, guid primitive.ObjectID) (*models.Group, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"group_id": guid,
		"members": bson.M{
			"$elemMatch": bson.M{"$eq": user_id},
		},
	}
	var result models.Group

	err := initializers.GroupCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateGroup(group *models.Group) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := initializers.GroupCollection.InsertOne(ctx, group)
	if err != nil {
		return err
	}

	return nil
}

