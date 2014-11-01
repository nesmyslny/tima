package server

import (
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"
)

type UserApi struct {
	db   *Db
	auth *Auth
}

const bcryptCost int = 13

func NewUserApi(db *Db, auth *Auth) *UserApi {
	return &UserApi{db, auth}
}

func (this *UserApi) SigninHandler(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	var credentials UserCredentials
	err := unmarshalJson(r.Body, &credentials)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}

	token, err := this.auth.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return nil, &CtrlHandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	return jsonResultString(token)
}

func (this *UserApi) IsSignedInHandler(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	signedIn := this.auth.ValidateToken(r)
	return jsonResultBool(signedIn)
}

func (this *UserApi) AddUser(username string, pwd string, firstName string, lastName string, email string) (*User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &User{
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
