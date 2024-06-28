package database

import (
	"context"
	"fmt"
	"study-bean/initializers"
	"study-bean/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindUserFromDatabase(email string) (*models.User, error) {
	// create a filter to search for the email
	filter := bson.M{"email": email}

	// retrieving the first document that matches the filter
	var result models.User
	// check for errors in the finding
	err := initializers.UserCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	return &result, nil
}


func AddUserToDatabase(user models.User) error {

	newUser := models.User{
		Email: user.Email,
		Username: user.Username,
		Password: user.Password,
	}

	// insert the bson object using InsertOne()
	_, err := initializers.UserCollection.InsertOne(context.TODO(), &newUser)
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