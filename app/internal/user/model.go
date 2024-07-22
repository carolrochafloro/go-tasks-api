package user

import (
	"errors"
	"go-tasks-api/app/internal/tasks"
	"strings"
)

type User struct {
	Id        string `bson:"_id,omitempty"`
	Email     string `bson:"email,omitempty"`
	Password  string `bson:"password,omitempty"`
	FirstName string `bson:"first_name,omitempty"`
	LastName  string `bson:"last_name,omitempty"`
	Tasks []tasks.Task `bson:"tasks,omitempty"`    
}

func (u *User) Validate() error {

	if !strings.Contains(u.Email, "@") || len(u.Email) < 4 {
		return errors.New("this email is not valid")
	}

	if len(u.FirstName) < 3 || len(u.LastName) < 3 {
		return errors.New("name and last name must be at least 3 characters long")
	}

	return nil
}


