package server

import "time"

type User struct {
	Id           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	PasswordHash []byte `db:"password_hash" json:"-"`
	FirstName    string `db:"first_name" json:"firstName"`
	LastName     string `db:"last_name" json:"lastName"`
	Email        string `db:"email" json:"email"`
}

type Project struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type Activity struct {
	Id           int       `db:"id" json:"id"`
	Day          time.Time `db:"day" json:"day"`
	UserId       int       `db:"user_id" json:"userId"`
	ProjectId    int       `db:"project_id" json:"projectId"`
	Duration     int       `db:"duration" json:"duration"`
	ProjectTitle string    `db:"-" json:"projectTitle"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonResult struct {
	BoolResult   bool   `json:"boolResult"`
	StringResult string `json:"stringResult"`
	IntResult    int    `json:"intResult"`
}
