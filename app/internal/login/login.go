package login

import (
	"context"
	"encoding/json"
	"go-tasks-api/app/internal/db"
	"go-tasks-api/app/internal/user"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()

	client := db.Client
	
	if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Content type is not application/json", http.StatusUnsupportedMediaType)
        return
    }

	if r.Body == nil {
        http.Error(w, "Request body is missing", http.StatusBadRequest)
        return
    }

	defer r.Body.Close()

	login := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	collection := client.Database("go-tasks").Collection("users")

	email := login["email"]
	password := login["password"]

	filter := bson.D{{Key: "email", Value: email}}

	var user user.UserT

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
            return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.Id,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful login"))

}