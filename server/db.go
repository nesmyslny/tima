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
	dbAccess.dbMap.AddTableWithName(User{}, "user").SetKeys(true, "id")
	dbAccess.dbMap.AddTableWithName(Project{}, "project").SetKeys(true, "id")
	dbAccess.dbMap.AddTableWithName(ProjectCategory{}, "project_category").SetKeys(true, "id")
	dbAccess.dbMap.AddTableWithName(ActivityType{}, "activity_type").SetKeys(true, "id")
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
	var err error
	if user.ID < 0 {
		err = db.dbMap.Insert(user)
	} else {
		// password hash is only provided, if password needs to be changed. to prevent empty password hashes, this
		// statement is explicitly specified.
		sql := "update user set username = ?, first_name = ?, last_name = ?, email = ?%s where id = ?"

		if len(user.PasswordHash) == 0 {
			sql := fmt.Sprintf(sql, "")
			_, err = db.dbMap.Exec(sql, user.Username, user.FirstName, user.LastName, user.Email, user.ID)
		} else {
			sql := fmt.Sprintf(sql, ", password_hash = ?")
			_, err = db.dbMap.Exec(sql, user.Username, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.ID)
		}
	}
	return err
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
	var err error
	if activity.ID < 0 {
		err = db.dbMap.Insert(activity)
	} else {
		_, err = db.dbMap.Update(activity)
	}
	return err
}

func (db *DB) TryGetActivity(day time.Time, userID int, projectID int, activityTypeID int) (*Activity, error) {
	var activity *Activity
	err := db.dbMap.SelectOne(&activity,
		"select * from activity where user_id = ? and day = ? and project_id = ? and activity_type_id = ?",
		userID, day.Format(dateLayout), projectID, activityTypeID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return activity, nil
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
	project.ActivityTypes, err = db.getActivityTypes(project.ID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (db *DB) GetProjects() ([]Project, error) {
	var projects []Project
	_, err := db.dbMap.Select(&projects, "select * from project order by ref_id_complete")
	if err != nil {
		return nil, err
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

func (db *DB) SaveProject(project *Project, addedActivityTypes []ProjectActivityType, removedActivityTypes []ProjectActivityType) error {
	trans, err := db.dbMap.Begin()
	if err != nil {
		return err
	}

	project.RefIDComplete, err = db.getProjectRefIDComplete(trans, project)
	if err != nil {
		return err
	}

	if project.ID < 0 {
		err = trans.Insert(project)
	} else {
		_, err = trans.Update(project)
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

func (db *DB) getActivityTypes(projectID int) ([]ActivityType, error) {
	var activityTypes []ActivityType
	_, err := db.dbMap.Select(&activityTypes, "select * from activity_type where id in (select activity_type_id from project_activity_type where project_id = ?)", projectID)
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (db *DB) GetProjectActivityTypes(projectID int) ([]ProjectActivityType, error) {
	var projectActivityTypes []ProjectActivityType
	_, err := db.dbMap.Select(&projectActivityTypes, "select * from project_activity_type where project_id = ?", projectID)
	if err != nil {
		return nil, err
	}
	return projectActivityTypes, nil
}

func (db *DB) SaveActivityType(activityType *ActivityType) error {
	var err error
	if activityType.ID < 0 {
		err = db.dbMap.Insert(activityType)
	} else {
		_, err = db.dbMap.Update(activityType)
	}
	return err
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

func (db *DB) GetProjectActivityTypeViewList() ([]ProjectActivityTypeView, error) {
	sql := "select pat.*, p.ref_id_complete project_ref_id_complete, p.title project_title, at.title activity_type_title " +
		"from project_activity_type pat, project p, activity_type at " +
		"where pat.project_id = p.id and pat.activity_type_id = at.id " +
		"order by p.title, at.title"

	var list []ProjectActivityTypeView
	_, err := db.dbMap.Select(&list, sql)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (db *DB) getProjectCategories(parent *ProjectCategory) ([]ProjectCategory, error) {
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

func (db *DB) GetProjectCategoryTree(parent *ProjectCategory) ([]ProjectCategory, error) {
	projectCategories, err := db.getProjectCategories(parent)
	if err != nil {
		return nil, err
	}

	for i := range projectCategories {
		projectCategories[i].ProjectCategories, err = db.GetProjectCategoryTree(&projectCategories[i])
		if err != nil {
			return nil, err
		}
	}

	return projectCategories, nil
}

func (db *DB) GetProjectCategoryList(parent *ProjectCategory) ([]ProjectCategory, error) {
	projectCategories, err := db.getProjectCategories(parent)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(projectCategories); i++ {
		children, err := db.GetProjectCategoryList(&projectCategories[i])
		if err != nil {
			return nil, err
		}

		// inserting children into slice after parent
		slicingIndex := i + 1
		projectCategories = append(projectCategories[:slicingIndex], append(children, projectCategories[slicingIndex:]...)...)
		i += len(children)
	}

	return projectCategories, nil
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

func (db *DB) updateProjectCategoryRefIDComplete(trans *gorp.Transaction, projectCategory *ProjectCategory) error {
	categories, err := db.GetProjectCategoryList(projectCategory)
	if err != nil {
		return err
	}

	for i := range categories {
		categories[i].RefIDComplete, err = db.getProjectCategoryRefIDComplete(trans, &categories[i])
		if err != nil {
			return err
		}
		if _, err = trans.Update(&categories[i]); err != nil {
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
		_, err = trans.Update(&projects[i])
		if err != nil {
			return err
		}
	}

	for i := range projectCategory.ProjectCategories {
		db.updateProjectRefIDComplete(trans, &projectCategory.ProjectCategories[i])
	}

	return nil
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

	if projectCategory.ID < 0 {
		err = trans.Insert(projectCategory)
		if err != nil {
			trans.Rollback()
			return err
		}
	} else {
		_, err = trans.Update(projectCategory)
		if err != nil {
			trans.Rollback()
			return err
		}

		err = db.updateProjectCategoryRefIDComplete(trans, projectCategory)
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

	children, err := db.getProjectCategories(projectCategory)
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

func (db *DB) GetUsers() ([]User, error) {
	var users []User
	_, err := db.dbMap.Select(&users, "select * from user order by username")
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
