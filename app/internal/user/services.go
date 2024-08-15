package user

import (
	"context"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func addUserToDB(u UserT) error {

	client := db.Client

	collection := client.Database("go-tasks").Collection("users")

	result, err := collection.InsertOne(context.TODO(), u)

	if err != nil {
		logging.Warn("Unable to insert user.")
		println(result)
		return err
	}

	return nil
}

func deleteUserService(s string) (*mongo.DeleteResult, error) {
	client := db.Client

	collection := client.Database("go-tasks").Collection("users")

	objectID := convertStringToId(s)

	filter := bson.D{{Key: "_id", Value: objectID}}

	result, err := collection.DeleteOne(context.TODO(), filter)
	
		if (err != nil) {
			logging.Error("Unable to delete user.", err)
		}

	return result, nil

}

func convertStringToId(s string) (primitive.ObjectID) {

	objectID, err := primitive.ObjectIDFromHex(s)

	if err != nil {
		logging.Error("Error on converting string to ObjectId", err)
	}

	return objectID

}

func updateUser(s string, u UserT) (*mongo.UpdateResult) {

	client := db.Client

	objectID := convertStringToId(s)

	collection := client.Database("go-tasks").Collection("users")
	filter := bson.D{{Key: "_id", Value: objectID}}

	updateFields := bson.D{}

	if u.Email != "" {
		updateFields = append(updateFields, bson.E{Key: "email", Value: u.Email})
	}
	if u.Password != "" {
		updateFields = append(updateFields, bson.E{Key: "password", Value: u.Password})
	}
	if u.FirstName != "" {
		updateFields = append(updateFields, bson.E{Key: "first_name", Value: u.FirstName})
	}
	if u.LastName != "" {
		updateFields = append(updateFields, bson.E{Key: "last_name", Value: u.LastName})
	}

	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	
	if err != nil {
		logging.Error("There whas an error on updating the user", err)
	}

	return result

}
