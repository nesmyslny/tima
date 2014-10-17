package models

type User struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	Email        string `db:"email"`
}
