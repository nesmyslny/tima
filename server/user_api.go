package server

import "net/http"

type UserApi struct {
	db   *Db
	auth *Auth
}

func NewUserApi(db *Db, auth *Auth) *UserApi {
	return &UserApi{db, auth}
}

func (this *UserApi) SigninHandler(w http.ResponseWriter, r *http.Request) (interface{}, *HandlerError) {
	var credentials UserCredentials
	err := unmarshalJson(r.Body, &credentials)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	token, err := this.auth.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return nil, &HandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	return jsonResultString(token)
}

func (this *UserApi) IsSignedInHandler(w http.ResponseWriter, r *http.Request) (interface{}, *HandlerError) {
	signedIn := this.auth.ValidateToken(r)
	return jsonResultBool(signedIn)
}

func (this *UserApi) AddUser(username string, pwd string, firstName string, lastName string, email string) (*User, error) {
	pwdHash, err := this.auth.GeneratePasswordHash(pwd)
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
