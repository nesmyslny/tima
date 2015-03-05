package server

import (
	"errors"
	"net/http"
)

type UserAPI struct {
	db   *DB
	auth *Auth
}

func NewUserAPI(db *DB, auth *Auth) *UserAPI {
	return &UserAPI{db, auth}
}

func (userAPI *UserAPI) SigninHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var credentials UserCredentials
	err := context.GetReqBodyJSON(&credentials)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	user := userAPI.db.GetUserByName(credentials.Username)
	if user == nil {
		return nil, &HandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	token, err := userAPI.auth.Authenticate(user, credentials.Password)
	if err != nil {
		return nil, &HandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	return jsonResultString(token)
}

func (userAPI *UserAPI) IsSignedInHandler(context *HandlerContext) (interface{}, *HandlerError) {
	signedIn := userAPI.auth.ValidateToken(context)
	return jsonResultBool(signedIn)
}

func (userAPI *UserAPI) authorizeGetSave(requestUserId int, user *User) (bool, error) {
	return *user.Role == RoleAdmin || requestUserId == user.ID, nil
}

func (userAPI *UserAPI) AuthorizeGet(context *HandlerContext) (bool, error) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return false, err
	}
	return userAPI.authorizeGetSave(id, context.User)
}

func (userAPI *UserAPI) AuthorizeSave(context *HandlerContext) (bool, error) {
	var user User
	err := context.GetReqBodyJSON(&user)
	if err != nil {
		return false, err
	}
	return userAPI.authorizeGetSave(user.ID, context.User)
}

func (userAPI *UserAPI) GetHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	requestedUser, err := userAPI.get(id)
	if err != nil {
		return nil, &HandlerError{err, "Error: User could not be found.", http.StatusBadRequest}
	}
	return requestedUser, nil
}

func (userAPI *UserAPI) GetListHandler(context *HandlerContext) (interface{}, *HandlerError) {
	users, err := userAPI.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve users", http.StatusInternalServerError}
	}
	return users, nil
}

func (userAPI *UserAPI) SaveHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var user User
	err := context.GetReqBodyJSON(&user)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = userAPI.save(&user, *context.User.Role == RoleAdmin)
	if err != nil {
		if err == errUsernameUnavailable {
			return nil, &HandlerError{err, "Error: Specified Username is not available.", http.StatusBadRequest}
		}
		// return nil, &HandlerError{err, "Error: User could not be saved.", http.StatusInternalServerError}
		return nil, &HandlerError{err, err.Error(), http.StatusInternalServerError}
	}
	return jsonResultInt(user.ID)
}

func (userAPI *UserAPI) AddUser(username string, role int, departmentId int, pwd string, firstName string, lastName string, email string) (*User, error) {
	pwdHash, err := userAPI.auth.GeneratePasswordHash(pwd)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           -1,
		Role:         &role,
		DepartmentID: &departmentId,
		Username:     username,
		PasswordHash: pwdHash,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
	}

	err = userAPI.db.SaveUser(user, true)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userAPI *UserAPI) get(id int) (*User, error) {
	user, err := userAPI.db.GetUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userAPI *UserAPI) getList() ([]User, error) {
	users, err := userAPI.db.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userAPI *UserAPI) save(user *User, saveAsAdmin bool) error {
	var err error

	if user.ID < 0 {
		available, err := userAPI.db.IsUsernameAvailable(user.Username)
		if err != nil {
			return err
		}
		if !available {
			return errUsernameUnavailable
		}
	}

	if user.NewPassword != "" {
		if user.NewPassword == user.NewPasswordConfirm {
			user.PasswordHash, err = userAPI.auth.GeneratePasswordHash(user.NewPassword)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Passwords do not match")
		}
	}

	return userAPI.db.SaveUser(user, saveAsAdmin)
}
