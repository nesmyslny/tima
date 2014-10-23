package services

import (
	"code.google.com/p/go.crypto/bcrypt"
	"gnomon/dbaccess"
	"gnomon/models"
)

type UserService struct {
	db *DbAccess.Db
}

const bcryptCost int = 13

func NewUserService(db *DbAccess.Db) *UserService {
	return &UserService{db}
}

func (this *UserService) AddUser(username string, pwd string, firstName string, lastName string, email string) (*models.User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Id:           -1,
		Username:     username,
		PasswordHash: pwdHash,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
	}

	err = this.db.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
