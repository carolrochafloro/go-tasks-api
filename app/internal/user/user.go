package user

import (
	"errors"
	"go-tasks-api/app/internal/tasks"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserT struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string `bson:"email,omitempty"`
	Password  string `bson:"password,omitempty"`
	FirstName string `bson:"first_name,omitempty"`
	LastName  string `bson:"last_name,omitempty"`
	Tasks []tasks.TaskT `bson:"tasks,omitempty"`    
}

func (u *UserT) Validate() error {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(\.[a-zA-Z]{2})?$`
	result := regexp.MustCompile(emailRegex)

	if !result.MatchString(u.Email) {
		return errors.New("this email is not valid")
	}

	if len(u.FirstName) < 3 || len(u.LastName) < 3 {
		return errors.New("name and last name must be at least 3 characters long")
	}

	return nil
}




