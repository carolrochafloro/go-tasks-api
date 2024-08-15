package middleware

import (
	"go-tasks-api/app/internal/logging"
	"go-tasks-api/app/internal/user"
	"go-tasks-api/app/internal/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


func Authenticate(data map[string]string) (string, error) {

	godotenv.Load()

	email := data["email"]
	password := data["password"]

	var user user.UserT

	err := utils.GetByKey(email, "email", "users", &user)

	if err != nil {

		logging.Error("User not found.", err)
		return "User not found", err

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "Wrong password", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.Id,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		return "Unable to generate token", err
	}

	return tokenString, nil 

}