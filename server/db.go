package server

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // this package is only imported because of gorp. it's not directly used here.
	"github.com/rubenv/sql-migrate"
	"gopkg.in/gorp.v1"
)

type DB struct {
	dbMap            *gorp.DbMap
	connectionString string
	dialect          string
	migrationDir     string
	migrationTable   string
}

func NewDB(connectionString string) *DB {
	dbAccess := &DB{
		connectionString: connectionString,
		dialect:          "mysql",
		migrationDir:     "migration",
		migrationTable:   "migration",
	}

	db, err := sql.Open(dbAccess.dialect, dbAccess.connectionString)
	if err != nil {
		// todo: error handling
		panic(err.Error())
	}

	dbAccess.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbAccess.dbMap.AddTableWithName(Department{}, "department").SetKeys(true, "id").SetVersionCol("version")
	dbAccess.dbMap.AddTableWithName(User{}, "user").SetKeys(true, "id").SetVersionCol("version")
	dbAccess.dbMap.AddTableWithName(Project{}, "project").SetKeys(true, "id").SetVersionCol("version")
	dbAccess.dbMap.AddTableWithName(ProjectDepartment{}, "project_department").SetKeys(false, "project_id", "department_id")
	dbAccess.dbMap.AddTableWithName(ProjectUser{}, "project_user").SetKeys(false, "project_id", "user_id")
	dbAccess.dbMap.AddTableWithName(ProjectCategory{}, "project_category").SetKeys(true, "id").SetVersionCol("version")
	dbAccess.dbMap.AddTableWithName(ActivityType{}, "activity_type").SetKeys(true, "id").SetVersionCol("version")
	dbAccess.dbMap.AddTableWithName(ProjectActivityType{}, "project_activity_type").SetKeys(false, "project_id", "activity_type_id")
	dbAccess.dbMap.AddTableWithName(Activity{}, "activity").SetKeys(true, "id")

	return dbAccess
}

func (db *DB) Close() error {
	return db.dbMap.Db.Close()
}

func (db *DB) initMigrate() migrate.MigrationSource {
	migrate.SetTable(db.migrationTable)
	migrationSource := &migrate.FileMigrationSource{
		Dir: db.migrationDir,
	}
	return migrationSource
}

func (db *DB) Upgrade(max int) error {
	migrationSource := db.initMigrate()
	_, err := migrate.ExecMax(db.dbMap.Db, db.dialect, migrationSource, migrate.Up, max)
	if err != nil {
		// todo: logging
		// if(!) any migration were applied, try to roll back:
		// migrate.ExecMax(db.dbMap.Db, db.dialect, migrations, migrate.Down, applied)
		return err
	}

	return nil
}

func (db *DB) Downgrade(max int) error {
	migrationSource := db.initMigrate()
	_, err := migrate.ExecMax(db.dbMap.Db, db.dialect, migrationSource, migrate.Down, max)
	return err
}

func (db *DB) Update(trans *gorp.Transaction, model interface{}) error {
	var err error
	if trans != nil {
		_, err = trans.Update(model)
	} else {
		_, err = db.dbMap.Update(model)
	}

	if _, ok := err.(gorp.OptimisticLockError); ok {
		return errOptimisticLocking
	}
	return err
}

func (db *DB) GetNumberOfUsers() (int, error) {
	count, err := db.dbMap.SelectInt("select count(*) from user")
	return int(count), err
}

func (db *DB) GetUserByName(username string) *User {
	var user *User
	err := db.dbMap.SelectOne(&user, "select * from user where username = ?", username)
	if err != nil {
		return nil
	}
	return user
}

func (db *DB) SaveUser(user *User) error {
	if user.ID < 0 {
		return db.dbMap.Insert(user)
	}
	return db.Update(nil, user)
}

func (db *DB) GetActivitiesByDay(userID int, day time.Time) ([]ActivityView, error) {
	sql := "select a.*, p.title project_title, at.title activity_type_title " +
		"from activity a, project p, activity_type at " +
		"where a.project_id = p.id and a.activity_type_id = at.id " +
		"and user_id = ? and day = ? " +
		"order by duration desc"

	var activities []ActivityView
	_, err := db.dbMap.Select(&activities, sql, userID, day.Format(dateLayout))
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (db *DB) GetActivity(id int) (*Activity, error) {
	obj, err := db.dbMap.Get(Activity{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Activity), nil
}

func (db *DB) SaveActivity(activity *Activity) error {
	if activity.ID < 0 {
		return db.dbMap.Insert(activity)
	}
	return db.Update(nil, activity)
}

func (db *DB) IsProjectReferenced(id int) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from activity where project_id = ?)", id)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (db *DB) DeleteActivity(activity *Activity) error {
	_, err := db.dbMap.Delete(activity)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetProject(id int) (*Project, error) {
	obj, err := db.dbMap.Get(Project{}, id)
	if err != nil {
		return nil, err
	}
	project := obj.(*Project)

	project.ActivityTypes, err = db.getActivityTypesOfProject(project.ID)
	if err != nil {
		return nil, err
	}

	project.Departments, err = db.getDepartmentsOfProject(project.ID)
	if err != nil {
		return nil, err
	}

	project.Users, err = db.getUsersOfProject(project.ID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (db *DB) GetProjects(departmentID *int) ([]Project, error) {
	sqlTemplate := "select * from project %sorder by ref_id_complete"
	sqlWhere := ""
	var args []interface{}

	// if a department id is present, only select project of that (or a descendant) department.
	if departmentID != nil {
		deptIDs, err := db.GetDepartmentIDsDownward(*departmentID)
		if err != nil {
			return nil, err
		}

		deptArgString := createSqlArgString(len(deptIDs))
		args = sliceIntToInterface(deptIDs)
		sqlWhere = "where responsible_user_id in (select id from user where department_id in (" + deptArgString + ")) "
	}

	sql := fmt.Sprintf(sqlTemplate, sqlWhere)
	var projects []Project
	_, err := db.dbMap.Select(&projects, sql, args...)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (db *DB) GetProjectsOfUser(userID int) ([]Project, error) {
	var projects []Project
	_, err := db.dbMap.Select(&projects,
		"select * from project where responsible_user_id = ? or manager_user_id = ? "+
			"order by ref_id_complete", userID, userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (db *DB) GetProjectsForSelection(userID int, departmentID int) ([]Project, error) {
	deptIDs, err := db.GetDepartmentIDsUpward(departmentID)
	if err != nil {
		return nil, err
	}
	deptArgString := createSqlArgString(len(deptIDs))
	deptIDsInterface := sliceIntToInterface(deptIDs)

	sql := "select * from project where " +
		"responsible_user_id = ? or manager_user_id = ? or id in (select project_id from project_user where user_id = ?) " +
		"or id in (select project_id from project_department where department_id in (" + deptArgString + ")) " +
		"order by title"

	var projects []Project
	args := []interface{}{userID, userID, userID}
	args = append(args, deptIDsInterface...)

	_, err = db.dbMap.Select(&projects, sql, args...)
	if err != nil {
		return nil, err
	}

	for i, project := range projects {
		projects[i].ActivityTypes, err = db.getActivityTypesOfProject(project.ID)
		if err != nil {
			return nil, err
		}
	}

	return projects, nil
}

func (db *DB) getProjectsByProjectCategory(trans *gorp.Transaction, projectCategoryId int) ([]Project, error) {
	var projects []Project
	_, err := trans.Select(&projects, "select * from project where project_category_id = ? order by ref_id_complete", projectCategoryId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (db *DB) getProjectRefIDComplete(trans *gorp.Transaction, project *Project) (string, error) {
	category, err := db.getProjectCategory(trans, project.ProjectCategoryID)
	if err != nil {
		return "", err
	}

	refIDCategory, err := db.getProjectCategoryRefIDComplete(trans, category)
	if err != nil {
		return "", err
	}

	return refIDCategory + "/" + project.RefID, nil
}

func (db *DB) isProjectRefIDCompleteUnique(projectID int, refIDComplete string) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from project where id != ? and ref_id_complete = ?)", projectID, refIDComplete)
	if err != nil {
		return false, err
	}
	return exists == 0, nil
}

func (db *DB) SaveProject(project *Project,
	addedActivityTypes []ProjectActivityType, removedActivityTypes []ProjectActivityType,
	addedDepartments []ProjectDepartment, removedDepartments []ProjectDepartment,
	addedUsers []ProjectUser, removedUsers []ProjectUser) error {
	trans, err := db.dbMap.Begin()
	if err != nil {
		return err
	}

	project.RefIDComplete, err = db.getProjectRefIDComplete(trans, project)
	if err != nil {
		return err
	}

	isUnique, err := db.isProjectRefIDCompleteUnique(project.ID, project.RefIDComplete)
	if err != nil {
		return err
	}
	if !isUnique {
		return errIDNotUnique
	}

	if project.ID < 0 {
		err = trans.Insert(project)
	} else {
		err = db.Update(trans, project)
	}
	if err != nil {
		trans.Rollback()
		return err
	}

	for _, removedActivityType := range removedActivityTypes {
		count, err := trans.Delete(&removedActivityType)
		if err != nil || count != 1 {
			trans.Rollback()
			return err
		}
	}
	for _, addedActivityType := range addedActivityTypes {
		// project id must be set, if it's a new project
		if addedActivityType.ProjectID < 0 {
			addedActivityType.ProjectID = project.ID
		}
		err = trans.Insert(&addedActivityType)
		if err != nil {
			trans.Rollback()
			return err
		}
	}

	for _, removedDepartment := range removedDepartments {
		count, err := trans.Delete(&removedDepartment)
		if err != nil || count != 1 {
			trans.Rollback()
			return err
		}
	}
	for _, addedDepartment := range addedDepartments {
		if addedDepartment.ProjectID < 0 {
			addedDepartment.ProjectID = project.ID
		}
		err = trans.Insert(&addedDepartment)
		if err != nil {
			trans.Rollback()
			return err
		}
	}

	for _, removedUser := range removedUsers {
		count, err := trans.Delete(&removedUser)
		if err != nil || count != 1 {
			trans.Rollback()
			return err
		}
	}
	for _, addedUser := range addedUsers {
		if addedUser.ProjectID < 0 {
			addedUser.ProjectID = project.ID
		}
		err = trans.Insert(&addedUser)
		if err != nil {
			trans.Rollback()
			return err
		}
	}

	return trans.Commit()
}

func (db *DB) DeleteProject(project *Project) error {
	_, err := db.dbMap.Delete(project)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetActivityType(id int) (*ActivityType, error) {
	obj, err := db.dbMap.Get(ActivityType{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*ActivityType), nil
}

func (db *DB) GetActivityTypes() ([]ActivityType, error) {
	var activityTypes []ActivityType
	_, err := db.dbMap.Select(&activityTypes, "select * from activity_type order by title")
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (db *DB) getActivityTypesOfProject(projectID int) ([]ActivityType, error) {
	var activityTypes []ActivityType
	_, err := db.dbMap.Select(&activityTypes, "select * from activity_type where id in (select activity_type_id from project_activity_type where project_id = ?)", projectID)
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (db *DB) GetProjectActivityTypeIDs(projectID int) ([]int, error) {
	var IDs []int
	_, err := db.dbMap.Select(&IDs, "select activity_type_id from project_activity_type where project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

func (db *DB) SaveActivityType(activityType *ActivityType) error {
	if activityType.ID < 0 {
		return db.dbMap.Insert(activityType)
	}
	return db.Update(nil, activityType)
}

func (db *DB) DeleteActivityType(activityType *ActivityType) error {
	_, err := db.dbMap.Delete(activityType)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsActivityTypeReferenced(activityTypeID int, projectID *int) (bool, error) {
	var exists int64
	var err error

	if projectID == nil {
		exists, err = db.dbMap.SelectInt("select exists(select id from activity where activity_type_id = ?)", activityTypeID)
	} else {
		exists, err = db.dbMap.SelectInt("select exists(select id from activity where activity_type_id = ? and project_id = ?)", activityTypeID, projectID)
	}

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (db *DB) GetDepartmentIDsUpward(startDeptID int) ([]int, error) {
	deptIDs := []int{startDeptID}

	parentDeptID, err := db.dbMap.SelectNullInt("select parent_id from department where id = ?", startDeptID)
	if err != nil {
		return nil, err
	}

	if parentDeptID.Valid {
		parentDeptIDs, err := db.GetDepartmentIDsUpward(int(parentDeptID.Int64))
		if err != nil {
			return nil, err
		}
		deptIDs = append(deptIDs, parentDeptIDs...)
	}

	return deptIDs, nil
}

func (db *DB) GetDepartmentIDsDownward(startDeptID int) ([]int, error) {
	deptIDs := []int{startDeptID}

	var childrenDeptIDs []int
	_, err := db.dbMap.Select(&childrenDeptIDs, "select id from department where parent_id = ?", startDeptID)
	if err != nil {
		return nil, err
	}

	for _, id := range childrenDeptIDs {
		ids, err := db.GetDepartmentIDsDownward(id)
		if err != nil {
			return nil, err
		}
		deptIDs = append(deptIDs, ids...)
	}

	return deptIDs, nil
}

func (db *DB) GetProjectCategories(parent *ProjectCategory) ([]ProjectCategory, error) {
	var projectCategories []ProjectCategory
	const sqlTemplate string = "select * from project_category where parent_id %s order by ref_id_complete"
	var err error

	if parent == nil {
		sql := fmt.Sprintf(sqlTemplate, "is null")
		_, err = db.dbMap.Select(&projectCategories, sql)
	} else {
		sql := fmt.Sprintf(sqlTemplate, "= ?")
		_, err = db.dbMap.Select(&projectCategories, sql, parent.ID)
	}
	if err != nil {
		return nil, err
	}

	for i := range projectCategories {
		db.setProjectCategoryPath(&projectCategories[i], parent)
	}

	return projectCategories, nil
}

func (db *DB) setProjectCategoryPath(projectCategory *ProjectCategory, parentCategory *ProjectCategory) {
	parentPath := ""
	if parentCategory != nil {
		parentPath = parentCategory.Path + " \u203A "
	}
	projectCategory.Path = parentPath + projectCategory.Title
}

func (db *DB) getProjectCategory(trans *gorp.Transaction, id int) (*ProjectCategory, error) {
	var obj interface{}
	var err error

	if trans == nil {
		obj, err = db.dbMap.Get(ProjectCategory{}, id)
	} else {
		obj, err = trans.Get(ProjectCategory{}, id)
	}

	if err != nil {
		return nil, err
	}
	return obj.(*ProjectCategory), nil
}

func (db *DB) GetProjectCategory(id int) (*ProjectCategory, error) {
	return db.getProjectCategory(nil, id)
}

func (db *DB) getProjectCategoryRefIDComplete(trans *gorp.Transaction, projectCategory *ProjectCategory) (string, error) {
	refID := projectCategory.RefID

	if projectCategory.ParentID != nil {
		parent, err := db.getProjectCategory(trans, *projectCategory.ParentID)
		if err != nil {
			return "", err
		}
		parentRefIDComplete, err := db.getProjectCategoryRefIDComplete(trans, parent)
		if err != nil {
			return "", err
		}
		refID = parentRefIDComplete + refID
	}

	return refID, nil
}

func (db *DB) updateProjectCategoryRefIDComplete(trans *gorp.Transaction, projectCategories []ProjectCategory) error {
	var err error

	for i := range projectCategories {
		projectCategories[i].RefIDComplete, err = db.getProjectCategoryRefIDComplete(trans, &projectCategories[i])
		if err != nil {
			return err
		}
		if err = db.Update(trans, &projectCategories[i]); err != nil {
			return err
		}
		err = db.updateProjectCategoryRefIDComplete(trans, projectCategories[i].ProjectCategories)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) updateProjectRefIDComplete(trans *gorp.Transaction, projectCategory *ProjectCategory) error {
	projects, err := db.getProjectsByProjectCategory(trans, projectCategory.ID)
	if err != nil {
		return err
	}

	for i := range projects {
		projects[i].RefIDComplete, err = db.getProjectRefIDComplete(trans, &projects[i])
		if err != nil {
			return err
		}
		err = db.Update(trans, &projects[i])
		if err != nil {
			return err
		}
	}

	for i := range projectCategory.ProjectCategories {
		db.updateProjectRefIDComplete(trans, &projectCategory.ProjectCategories[i])
	}

	return nil
}

func (db *DB) isProjectCategoryRefIDCompleteUnique(projectCategoryID int, refIDComplete string) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from project_category where id != ? and ref_id_complete = ?)", projectCategoryID, refIDComplete)
	if err != nil {
		return false, err
	}
	return exists == 0, nil
}

func (db *DB) SaveProjectCategory(projectCategory *ProjectCategory) error {
	trans, err := db.dbMap.Begin()
	if err != nil {
		return err
	}

	projectCategory.RefIDComplete, err = db.getProjectCategoryRefIDComplete(trans, projectCategory)
	if err != nil {
		return err
	}

	isUnique, err := db.isProjectCategoryRefIDCompleteUnique(projectCategory.ID, projectCategory.RefIDComplete)
	if err != nil {
		return err
	}
	if !isUnique {
		return errIDNotUnique
	}

	if projectCategory.ID < 0 {
		err = trans.Insert(projectCategory)
		if err != nil {
			trans.Rollback()
			return err
		}
	} else {
		err = db.Update(trans, projectCategory)
		if err != nil {
			trans.Rollback()
			return err
		}

		err = db.updateProjectCategoryRefIDComplete(trans, projectCategory.ProjectCategories)
		if err != nil {
			trans.Rollback()
			return err
		}

		err = db.updateProjectRefIDComplete(trans, projectCategory)
		if err != nil {
			trans.Rollback()
			return err
		}
	}

	return trans.Commit()
}

func (db *DB) DeleteProjectCategory(projectCategory *ProjectCategory) error {
	_, err := db.dbMap.Delete(projectCategory)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsProjectCategoryReferenced(projectCategory *ProjectCategory) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from project where project_category_id = ?)", projectCategory.ID)
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	}

	children, err := db.GetProjectCategories(projectCategory)
	for _, child := range children {
		isReferenced, err := db.IsProjectCategoryReferenced(&child)
		if err != nil {
			return false, err
		} else if isReferenced {
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) GetUsers(departmentID *int) ([]User, error) {
	sqlTemplate := "select * from user %sorder by username"
	sqlWhere := ""
	var args []interface{}

	if departmentID != nil {
		deptIDs, err := db.GetDepartmentIDsDownward(*departmentID)
		if err != nil {
			return nil, err
		}

		deptArgString := createSqlArgString(len(deptIDs))
		args = sliceIntToInterface(deptIDs)
		sqlWhere = "where department_id in (" + deptArgString + ")"
	}

	sql := fmt.Sprintf(sqlTemplate, sqlWhere)
	var users []User
	_, err := db.dbMap.Select(&users, sql, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DB) getUsersOfProject(projectID int) ([]User, error) {
	var users []User
	_, err := db.dbMap.Select(&users, "select * from user where id in (select user_id from project_user where project_id = ?)", projectID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DB) GetUser(id int) (*User, error) {
	obj, err := db.dbMap.Get(User{}, id)
	if err != nil {
		return nil, err
	}

	return obj.(*User), nil
}

func (db *DB) IsUsernameAvailable(username string) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from user where username = ?)", username)
	if err != nil {
		return false, err
	}

	return exists == 0, nil
}

func (db *DB) GetDepartments(parent *Department) ([]Department, error) {
	var departments []Department
	const sqlTemplate string = "select * from department where parent_id %s order by title"
	var err error

	if parent == nil {
		sql := fmt.Sprintf(sqlTemplate, "is null")
		_, err = db.dbMap.Select(&departments, sql)
	} else {
		sql := fmt.Sprintf(sqlTemplate, "= ?")
		_, err = db.dbMap.Select(&departments, sql, parent.ID)
	}
	if err != nil {
		return nil, err
	}

	for i := range departments {
		db.setDepartmentPath(&departments[i], parent)
	}

	return departments, nil
}

func (db *DB) setDepartmentPath(department *Department, parentDepartment *Department) {
	parentPath := ""
	if parentDepartment != nil {
		parentPath = parentDepartment.Path + " \u203A "
	}
	department.Path = parentPath + department.Title
}

func (db *DB) GetDepartment(id int) (*Department, error) {
	obj, err := db.dbMap.Get(Department{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Department), nil
}

func (db *DB) SaveDepartment(department *Department) error {
	if department.ID < 0 {
		return db.dbMap.Insert(department)
	}
	return db.Update(nil, department)
}

func (db *DB) DeleteDepartment(department *Department) error {
	_, err := db.dbMap.Delete(department)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsDepartmentReferenced(department *Department) (bool, error) {
	exists, err := db.dbMap.SelectInt("select exists(select id from user where department_id = ?)", department.ID)
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	}

	children, err := db.GetDepartments(department)
	for _, child := range children {
		isReferenced, err := db.IsDepartmentReferenced(&child)
		if err != nil {
			return false, err
		} else if isReferenced {
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) GetProjectUserIDs(projectID int) ([]int, error) {
	var IDs []int
	_, err := db.dbMap.Select(&IDs, "select user_id from project_user where project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}

func (db *DB) getDepartmentsOfProject(projectID int) ([]Department, error) {
	var depts []Department
	_, err := db.dbMap.Select(&depts, "select * from department where id in (select department_id from project_department where project_id = ?)", projectID)
	if err != nil {
		return nil, err
	}
	return depts, nil
}

func (db *DB) GetProjectDepartmentIDs(projectID int) ([]int, error) {
	var IDs []int
	_, err := db.dbMap.Select(&IDs, "select department_id from project_department where project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	return IDs, nil
}
