package server

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql" // this package is only imported because of gorp. it's not directly used here.
	"github.com/rubenv/sql-migrate"
)

type DB struct {
	dbMap            *gorp.DbMap
	connectionString string
	dialect          string
	migrationDir     string
	migrationTable   string
}

const dateLayout string = "2006-01-02"

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
	dbAccess.dbMap.AddTableWithName(ActivityType{}, "activity_type").SetKeys(true, "id")
	dbAccess.dbMap.AddTableWithName(ProjectActivityType{}, "project_activity_type").SetKeys(false, "project_id", "activity_type_id")
	dbAccess.dbMap.AddTableWithName(Activity{}, "activity").SetKeys(true, "id")

	return dbAccess
}

func (db *DB) Close() error {
	return db.dbMap.Db.Close()
}

func (db *DB) Upgrade() error {
	migrate.SetTable(db.migrationTable)
	migrationSource := &migrate.FileMigrationSource{
		Dir: db.migrationDir,
	}

	_, err := migrate.Exec(db.dbMap.Db, db.dialect, migrationSource, migrate.Up)
	if err != nil {
		// todo: logging
		// if(!) any migration were applied, try to roll back:
		// migrate.ExecMax(db.dbMap.Db, db.dialect, migrations, migrate.Down, applied)
		return err
	}

	return nil
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
		_, err = db.dbMap.Update(user)
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
	_, err := db.dbMap.Select(&projects, "select * from project order by title")
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (db *DB) SaveProject(project *Project, addedActivityTypes []ProjectActivityType, removedActivityTypes []ProjectActivityType) error {
	trans, err := db.dbMap.Begin()
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
	sql := "select pat.*, p.title project_title, at.title activity_type_title " +
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

func (db *DB) GetProjectCategories(parentID *int) ([]ProjectCategory, error) {
	var projectCategories []ProjectCategory
	const sqlTemplate string = "select * from project_category where parent_id %s order by title"
	var err error

	if parentID == nil {
		sql := fmt.Sprintf(sqlTemplate, "is null")
		_, err = db.dbMap.Select(&projectCategories, sql)
	} else {
		sql := fmt.Sprintf(sqlTemplate, "= ?")
		_, err = db.dbMap.Select(&projectCategories, sql, *parentID)
	}
	if err != nil {
		return nil, err
	}

	for i, cat := range projectCategories {
		projectCategories[i].ProjectCategories, err = db.GetProjectCategories(&cat.ID)
		if err != nil {
			return nil, err
		}
	}

	return projectCategories, nil
}
