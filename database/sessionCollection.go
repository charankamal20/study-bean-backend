package database

import (
	"context"
	"errors"
	"study-bean/initializers"
	"study-bean/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewSession(session *models.Session) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := initializers.SessionCollection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

func JoinSession(uid string, session_id string, isLoggedIn bool) (error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{
		"session_id": session_id,
	}

	var update primitive.M
	if isLoggedIn {
		update = bson.M{
			"$addToSet": bson.M{
				"members": uid,
			},
				"$inc": bson.M{
				"number_of_members": +1,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}
	} else {
		update = bson.M{
			"$addToSet": bson.M{
				"temp_member_list": uid,
			},
			"$inc": bson.M{
				"number_of_members": +1,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}
	}

	result, err := initializers.SessionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no session found")
	}

	return nil
}