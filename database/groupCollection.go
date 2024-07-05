package database

import (
	"context"
	"study-bean/models"
	"time"
)

func AddUserToGroup(user_id string) error {

	return nil
}

func FindUserInGroup(user_id string) (*models.User, error) {

	_, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	// Get user info from database

	return nil, nil
}