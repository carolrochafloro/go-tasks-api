package user

import (
	"context"
	"go-tasks-api/app/internal/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUser(s, key string, client *mongo.Client) (UserT, bool) {

	if client == nil {
		logging.Error("MongoDB client is nil")
		return UserT{}, false
	}

	collection := client.Database("go-tasks").Collection("users")
	println("Test", s)
	
	var filter interface{}

	// converter id para ObjectId caso a busca seja por OId
	if key == "_id" {

		objectID, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			logging.Error("Invalid ObjectID format", err)
			return UserT{}, false
		}
		filter = bson.M{key: objectID}
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

func addUserToDB(u UserT, client *mongo.Client) {

	collection := client.Database("go-tasks").Collection("users")

	result, err := collection.InsertOne(context.TODO(), u)

	if err != nil {
		logging.Warn("Unable to insert user.")
		println(result)
		return
	}
}