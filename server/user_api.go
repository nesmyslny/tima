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

	return &SingleValue{token}, nil
}

func (userAPI *UserAPI) IsSignedInHandler(context *HandlerContext) (interface{}, *HandlerError) {
	signedIn := userAPI.auth.ValidateToken(context)
	return &SingleValue{signedIn}, nil
}

func (userAPI *UserAPI) authorizeGetSave(requestUserId int, user *User) bool {
	return *user.Role == RoleAdmin || requestUserId == user.ID
}

func (userAPI *UserAPI) AuthorizeGet(context *HandlerContext) (bool, error) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return false, err
	}
	return userAPI.authorizeGetSave(id, context.User), nil
}

func (userAPI *UserAPI) AuthorizeSave(context *HandlerContext) (bool, error) {
	var user User
	err := context.GetReqBodyJSON(&user)
	if err != nil {
		return false, err
	}
	return userAPI.authorizeGetSave(user.ID, context.User), nil
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
	return userAPI.getListHandler(nil)
}

func (userAPI *UserAPI) GetListDeptHandler(context *HandlerContext) (interface{}, *HandlerError) {
	return userAPI.getListHandler(context.User.DepartmentID)
}

func (userAPI *UserAPI) getListHandler(deptID *int) (interface{}, *HandlerError) {
	users, err := userAPI.getList(deptID)
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
		} else if err == errOptimisticLocking {
			return nil, &HandlerError{err, "Error: User was changed/deleted by another user.", http.StatusInternalServerError}
		}
		return nil, &HandlerError{err, "Error: User could not be saved.", http.StatusInternalServerError}
	}
	return user, nil
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

	err = userAPI.save(user, true)
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

func (userAPI *UserAPI) getList(departmentID *int) ([]User, error) {
	users, err := userAPI.db.GetUsers(departmentID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userAPI *UserAPI) save(user *User, saveAsAdmin bool) error {
	var passwordHashOrig []byte

	if user.ID < 0 {
		available, err := userAPI.db.IsUsernameAvailable(user.Username)
		if err != nil {
			return err
		}
		if !available {
			return errUsernameUnavailable
		}
	} else {
		userOrig, err := userAPI.db.GetUser(user.ID)
		if err != nil {
			return err
		}

		passwordHashOrig = userOrig.PasswordHash

		// some attributes may only be changed by admins -> reset these attributes in other cases.
		if !saveAsAdmin {
			user.Role = userOrig.Role
			user.DepartmentID = userOrig.DepartmentID
		}
	}

	err := userAPI.setPasswordBeforeSave(user, passwordHashOrig)
	if err != nil {
		return err
	}

	return userAPI.db.SaveUser(user)
}

func (userAPI *UserAPI) setPasswordBeforeSave(user *User, passwordHasOrig []byte) error {
	if user.NewPassword != "" {
		if user.NewPassword == user.NewPasswordConfirm {
			var err error
			user.PasswordHash, err = userAPI.auth.GeneratePasswordHash(user.NewPassword)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Passwords do not match")
		}
	} else {
		// password is only provided, if password needs to be changed. if it's not set, reset it to the original.
		user.PasswordHash = passwordHasOrig
	}
	return nil
}
