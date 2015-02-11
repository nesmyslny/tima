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
	Id            int            `db:"id" json:"id"`
	Title         string         `db:"title" json:"title"`
	ActivityTypes []ActivityType `db:"-" json:"activityTypes"`
}

type ActivityType struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type ProjectActivityTypes struct {
	ProjectId      int `db:"project_id"`
	ActivityTypeId int `db:"activity_type_id"`
}

type ProjectActivityTypesView struct {
	ProjectId         int    `db:"project_id" json:"projectId"`
	ActivityTypeId    int    `db:"activity_type_id" json:"activityTypeId"`
	ProjectTitle      string `db:"project_title" json:"projectTitle"`
	ActivityTypeTitle string `db:"activity_type_title" json:"activityTypeTitle"`
}

type Activity struct {
	Id             int       `db:"id" json:"id"`
	Day            time.Time `db:"day" json:"day"`
	UserId         int       `db:"user_id" json:"userId"`
	ProjectId      int       `db:"project_id" json:"projectId"`
	ActivityTypeId int       `db:"activity_type_id" json:"activityTypeId"`
	Duration       int       `db:"duration" json:"duration"`
}

type ActivityView struct {
	Id                int       `db:"id" json:"id"`
	Day               time.Time `db:"day" json:"day"`
	UserId            int       `db:"user_id" json:"userId"`
	ProjectId         int       `db:"project_id" json:"projectId"`
	ActivityTypeId    int       `db:"activity_type_id" json:"activityTypeId"`
	Duration          int       `db:"duration" json:"duration"`
	ProjectTitle      string    `db:"project_title" json:"projectTitle"`
	ActivityTypeTitle string    `db:"activity_type_title" json:"activityTypeTitle"`
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
