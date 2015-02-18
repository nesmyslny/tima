package server

import "time"

type User struct {
	ID           int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	PasswordHash []byte `db:"password_hash" json:"-"`
	FirstName    string `db:"first_name" json:"firstName"`
	LastName     string `db:"last_name" json:"lastName"`
	Email        string `db:"email" json:"email"`
}

type Project struct {
	ID                int            `db:"id" json:"id"`
	ProjectCategoryID int            `db:"project_category_id" json:"projectCategoryId"`
	Title             string         `db:"title" json:"title"`
	ActivityTypes     []ActivityType `db:"-" json:"activityTypes"`
}

type ProjectCategory struct {
	ID                int               `db:"id" json:"id"`
	ParentID          *int              `db:"parent_id" json:"parentId"`
	Title             string            `db:"title" json:"title"`
	Path              string            `db:"-" json:"path"`
	ProjectCategories []ProjectCategory `db:"-" json:"projectCategories"`
}

type ActivityType struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type ProjectActivityType struct {
	ProjectID      int `db:"project_id"`
	ActivityTypeID int `db:"activity_type_id"`
}

type ProjectActivityTypeView struct {
	ProjectID         int    `db:"project_id" json:"projectId"`
	ActivityTypeID    int    `db:"activity_type_id" json:"activityTypeId"`
	ProjectTitle      string `db:"project_title" json:"projectTitle"`
	ActivityTypeTitle string `db:"activity_type_title" json:"activityTypeTitle"`
}

type Activity struct {
	ID             int       `db:"id" json:"id"`
	Day            time.Time `db:"day" json:"day"`
	UserID         int       `db:"user_id" json:"userId"`
	ProjectID      int       `db:"project_id" json:"projectId"`
	ActivityTypeID int       `db:"activity_type_id" json:"activityTypeId"`
	Duration       int       `db:"duration" json:"duration"`
}

type ActivityView struct {
	ID                int       `db:"id" json:"id"`
	Day               time.Time `db:"day" json:"day"`
	UserID            int       `db:"user_id" json:"userId"`
	ProjectID         int       `db:"project_id" json:"projectId"`
	ActivityTypeID    int       `db:"activity_type_id" json:"activityTypeId"`
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
