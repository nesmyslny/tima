package server

import (
	"errors"
	"time"
)

var errItemInUse = errors.New("Item is already in use")
var errUsernameUnavailable = errors.New("Username unavailable")
var errIDNotUnique = errors.New("ID must be unique")
var errOptimisticLocking = errors.New("Data was changed/deleted")
var errForbidden = errors.New("Forbidden")

const RoleUser int = 10
const RoleDeptManager int = 30
const RoleManager int = 50
const RoleAdmin int = 99

const dateLayout string = "2006-01-02"
const sec8h int = 28800

type Department struct {
	ID          int          `db:"id" json:"id"`
	ParentID    *int         `db:"parent_id" json:"parentId"`
	Title       string       `db:"title" json:"title"`
	Version     int          `db:"version" json:"version"`
	Path        string       `db:"-" json:"path"`
	Departments []Department `db:"-" json:"departments"`
}

type User struct {
	ID                 int    `db:"id" json:"id"`
	Role               *int   `db:"role" json:"role"`
	DepartmentID       *int   `db:"department_id" json:"departmentId"`
	Username           string `db:"username" json:"username"`
	PasswordHash       []byte `db:"password_hash" json:"-"`
	FirstName          string `db:"first_name" json:"firstName"`
	LastName           string `db:"last_name" json:"lastName"`
	Email              string `db:"email" json:"email"`
	Version            int    `db:"version" json:"version"`
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
	Description       string         `db:"description" json:"description"`
	Version           int            `db:"version" json:"version"`
	ActivityTypes     []ActivityType `db:"-" json:"activityTypes"`
	Departments       []Department   `db:"-" json:"departments"`
	Users             []User         `db:"-" json:"users"`
}

func (project *Project) getActivityTypeIDs() []int {
	var IDs []int
	for _, activityType := range project.ActivityTypes {
		IDs = append(IDs, activityType.ID)
	}
	return IDs
}

func (project *Project) GetDepartmentIDs() []int {
	var IDs []int
	for _, dept := range project.Departments {
		IDs = append(IDs, dept.ID)
	}
	return IDs
}

func (project *Project) GetUserIDs() []int {
	var IDs []int
	for _, user := range project.Users {
		IDs = append(IDs, user.ID)
	}
	return IDs
}

type ProjectCategory struct {
	ID                int               `db:"id" json:"id"`
	ParentID          *int              `db:"parent_id" json:"parentId"`
	RefID             string            `db:"ref_id" json:"refId"`
	RefIDComplete     string            `db:"ref_id_complete" json:"refIdComplete"`
	Title             string            `db:"title" json:"title"`
	Version           int               `db:"version" json:"version"`
	Path              string            `db:"-" json:"path"`
	ProjectCategories []ProjectCategory `db:"-" json:"projectCategories"`
}

type ProjectDepartment struct {
	ProjectID    int `db:"project_id"`
	DepartmentID int `db:"department_id"`
}

type ProjectUser struct {
	ProjectID int `db:"project_id"`
	UserID    int `db:"user_id"`
}

type ActivityType struct {
	ID      int    `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Version int    `db:"version" json:"version"`
}

type ProjectActivityType struct {
	ProjectID      int `db:"project_id"`
	ActivityTypeID int `db:"activity_type_id"`
}

type Activity struct {
	ID             int       `db:"id" json:"id"`
	Day            time.Time `db:"day" json:"day"`
	UserID         int       `db:"user_id" json:"userId"`
	ProjectID      int       `db:"project_id" json:"projectId"`
	ActivityTypeID int       `db:"activity_type_id" json:"activityTypeId"`
	Duration       int       `db:"duration" json:"duration"`
	Description    string    `db:"description" json:"description"`
	Version        int       `db:"version" json:"version"`
}

type ActivityView struct {
	ID                int       `db:"id" json:"id"`
	Day               time.Time `db:"day" json:"day"`
	UserID            int       `db:"user_id" json:"userId"`
	ProjectID         int       `db:"project_id" json:"projectId"`
	ActivityTypeID    int       `db:"activity_type_id" json:"activityTypeId"`
	Duration          int       `db:"duration" json:"duration"`
	Description       string    `db:"description" json:"description"`
	Version           int       `db:"version" json:"version"`
	ProjectTitle      string    `db:"project_title" json:"projectTitle"`
	ActivityTypeTitle string    `db:"activity_type_title" json:"activityTypeTitle"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SingleValue struct {
	Value interface{} `json:"value"`
}
