package server

import "time"

const RoleUser int = 10
const RoleManager int = 30
const RoleAdmin int = 99

const dateLayout string = "2006-01-02"

type User struct {
	ID                 int    `db:"id" json:"id"`
	Role               *int   `db:"role" json:"role"`
	Username           string `db:"username" json:"username"`
	PasswordHash       []byte `db:"password_hash" json:"-"`
	FirstName          string `db:"first_name" json:"firstName"`
	LastName           string `db:"last_name" json:"lastName"`
	Email              string `db:"email" json:"email"`
	NewPassword        string `db:"-" json:"newPassword"`
	NewPasswordConfirm string `db:"-" json:"newPasswordConfirm"`
}

type Project struct {
	ID                int            `db:"id" json:"id"`
	ProjectCategoryID int            `db:"project_category_id" json:"projectCategoryId"`
	RefID             string         `db:"ref_id" json:"refId"`
	RefIDComplete     string         `db:"ref_id_complete" json:"refIdComplete"`
	ResponsibleUserID *int           `db:"responsible_user_id" json:"responsibleUserId"`
	ManagerUserID     *int           `db:"manager_user_id" json:"managerUserId"`
	Title             string         `db:"title" json:"title"`
	ActivityTypes     []ActivityType `db:"-" json:"activityTypes"`
}

type ProjectCategory struct {
	ID                int               `db:"id" json:"id"`
	ParentID          *int              `db:"parent_id" json:"parentId"`
	RefID             string            `db:"ref_id" json:"refId"`
	RefIDComplete     string            `db:"ref_id_complete" json:"refIdComplete"`
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
	ProjectID            int    `db:"project_id" json:"projectId"`
	ActivityTypeID       int    `db:"activity_type_id" json:"activityTypeId"`
	ProjectRefIDComplete string `db:"project_ref_id_complete" json:"projectRefIdComplete"`
	ProjectTitle         string `db:"project_title" json:"projectTitle"`
	ActivityTypeTitle    string `db:"activity_type_title" json:"activityTypeTitle"`
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
