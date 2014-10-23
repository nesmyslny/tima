package models

type User struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash []byte `db:"password_hash"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	Email        string `db:"email"`
}

type UserSignin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonResult struct {
	BoolResult   bool   `json:"BoolResult"`
	StringResult string `json:"StringResult"`
}
