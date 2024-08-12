package user

import (
	"context"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUser(s, key string) (UserT, bool) {

	client := db.Client

	if client == nil {
		logging.Error("MongoDB client is nil")
		return UserT{}, false
	}

	collection := client.Database("go-tasks").Collection("users")
	
	var filter interface{}

	// converter id para ObjectId caso a busca seja por OId
	if key == "_id" {

		object := convertStringToId(s)
		filter = bson.M{key: object}

	} else {

		filter = bson.M{key: s}
	}
	
	var user UserT
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return UserT{}, false
	}
	return user, true
}

func addUserToDB(u UserT) {

	client := db.Client

	collection := client.Database("go-tasks").Collection("users")

	result, err := collection.InsertOne(context.TODO(), u)

	if err != nil {
		logging.Warn("Unable to insert user.")
		println(result)
		return
	}
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