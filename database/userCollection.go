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

func AddGroupInUser(guid string, user_id string) error {

	filter := bson.M{"user_id": user_id}

	update := bson.M{
		"$addToSet": bson.M{
			"groups": guid,
		},
		"$set": bson.M{"updated_at": time.Now()},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.User
	err := initializers.UserCollection.FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

func FindUserByEmail(email string) (*models.User, error) {
	// create a filter to search for the email
	filter := bson.M{"email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// retrieving the first document that matches the filter
	var result models.User
	// check for errors in the finding
	err := initializers.UserCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &result, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	// create a filter to search for the username
	filter := bson.M{"username": username}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// retrieving the first document that matches the filter
	var result models.User
	// check for errors in the finding
	err := initializers.UserCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &result, nil
}

func FindUserByUserID(userID string) (*models.User, error) {
	// create a filter to search for the username
	filter := bson.M{"user_id": userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// retrieving the first document that matches the filter
	var result models.User
	// check for errors in the finding
	err := initializers.UserCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &result, nil
}

func AddUserToDatabase(user models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// insert the bson object using InsertOne()
	_, err := initializers.UserCollection.InsertOne(ctx, &user)
	fmt.Println(err)
	// check for errors in the insertion
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserRefreshToken(email string, refresh_token string) error {

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"refresh_token": refresh_token}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := initializers.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func AddTodoToUser(todo models.UserTodo) error {
	// insert the bson object using InsertOne()
	_, err := initializers.UserTodoCollection.InsertOne(context.TODO(), &todo)
	fmt.Println(err)
	// check for errors in the insertion
	if err != nil {
		return err
	}

	return nil
}

func FindAllUsers() ([]models.User, error) {

	usersCursor, err := initializers.UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var userList []models.User
	for usersCursor.Next(context.TODO()) {
		var user models.User
		err := usersCursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	return userList, nil
}

func UpdateUserByKey[T any](_id primitive.ObjectID, key string, newValue T) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := initializers.UserCollection.UpdateByID(ctx, _id, bson.M{key: newValue})
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExistInGroup(user_id string, guid string) (*models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"user_id": user_id,
		"groups": bson.M{
			"$elemMatch": bson.M{"$eq": guid},
		},
	}

	var user models.User

	err := initializers.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case when no documents match the filter
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
