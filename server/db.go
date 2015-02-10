package server

import (
	"database/sql"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rubenv/sql-migrate"
)

type Db struct {
	dbMap            *gorp.DbMap
	connectionString string
	dialect          string
	migrationDir     string
	migrationTable   string
}

const dateLayout string = "2006-01-02"

func NewDb(connectionString string) *Db {
	dbAccess := &Db{
		connectionString: connectionString,
		dialect:          "mysql",
		migrationDir:     "migrations",
		migrationTable:   "migrations",
	}

	db, err := sql.Open(dbAccess.dialect, dbAccess.connectionString)
	if err != nil {
		// todo: error handling
		panic(err.Error())
	}

	dbAccess.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbAccess.dbMap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbAccess.dbMap.AddTableWithName(Project{}, "projects").SetKeys(true, "Id")
	dbAccess.dbMap.AddTableWithName(ActivityType{}, "activity_types").SetKeys(true, "Id")
	dbAccess.dbMap.AddTableWithName(ProjectActivityTypes{}, "projects_activity_types").SetKeys(false, "project_id", "activity_type_id")
	dbAccess.dbMap.AddTableWithName(Activity{}, "activities").SetKeys(true, "Id")

	return dbAccess
}

func (this *Db) Close() error {
	return this.dbMap.Db.Close()
}

func (this *Db) Upgrade() error {
	migrate.SetTable(this.migrationTable)
	migrations := &migrate.FileMigrationSource{
		Dir: this.migrationDir,
	}

	_, err := migrate.Exec(this.dbMap.Db, this.dialect, migrations, migrate.Up)
	if err != nil {
		// todo: logging
		// if(!) any migration were applied, try to roll back:
		// migrate.ExecMax(this.dbMap.Db, this.dialect, migrations, migrate.Down, applied)
		return err
	}

	return nil
}

func (this *Db) GetNumberOfUsers() (int, error) {
	count, err := this.dbMap.SelectInt("select count(*) from users")
	return int(count), err
}

func (this *Db) GetUserByName(username string) *User {
	var user *User
	err := this.dbMap.SelectOne(&user, "select * from users where username = ?", username)
	if err != nil {
		return nil
	}
	return user
}

func (this *Db) SaveUser(user *User) error {
	var err error
	if user.Id < 0 {
		err = this.dbMap.Insert(user)
	} else {
		_, err = this.dbMap.Update(user)
	}
	return err
}

func (this *Db) GetActivitiesByDay(userId int, day time.Time) ([]Activity, error) {
	var activities []Activity
	_, err := this.dbMap.Select(&activities,
		"select * from activities where user_id = ? and day = ? order by duration desc",
		userId, day.Format(dateLayout))
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (this *Db) GetActivity(id int) (*Activity, error) {
	obj, err := this.dbMap.Get(Activity{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*Activity), nil
}

func (this *Db) SaveActivity(activity *Activity) error {
	var err error
	if activity.Id < 0 {
		err = this.dbMap.Insert(activity)
	} else {
		_, err = this.dbMap.Update(activity)
	}
	return err
}

func (this *Db) TryGetActivity(day time.Time, userId int, projectId int) (*Activity, error) {
	var activity *Activity
	err := this.dbMap.SelectOne(&activity,
		"select * from activities where user_id = ? and day = ? and project_id = ?",
		userId, day.Format(dateLayout), projectId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return activity, nil
}

func (this *Db) IsProjectReferenced(id int) (bool, error) {
	exists, err := this.dbMap.SelectInt("select exists(select id from activities where project_id = ?)", id)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (this *Db) DeleteActivity(activity *Activity) error {
	_, err := this.dbMap.Delete(activity)
	if err != nil {
		return err
	}

	return nil
}

func (this *Db) GetProject(id int) (*Project, error) {
	obj, err := this.dbMap.Get(Project{}, id)
	if err != nil {
		return nil, err
	}

	project := obj.(*Project)
	project.ActivityTypes, err = this.getProjectActivityTypes(project.Id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (this *Db) GetProjects() ([]Project, error) {
	var projects []Project
	_, err := this.dbMap.Select(&projects, "select * from projects order by title")
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (this *Db) SaveProject(project *Project) error {
	trans, err := this.dbMap.Begin()
	if err != nil {
		return err
	}

	err = this.deleteProjectActivityTypes(trans, project)
	if err != nil {
		trans.Rollback()
		return err
	}

	if project.Id < 0 {
		err = trans.Insert(project)
	} else {
		_, err = trans.Update(project)
	}
	return trans.Commit()
}

func (this *Db) DeleteProject(project *Project) error {
	_, err := this.dbMap.Delete(project)
	if err != nil {
		return err
	}

	return nil
}

func (this *Db) GetActivityType(id int) (*ActivityType, error) {
	obj, err := this.dbMap.Get(ActivityType{}, id)
	if err != nil {
		return nil, err
	}
	return obj.(*ActivityType), nil
}

func (this *Db) GetActivityTypes() ([]ActivityType, error) {
	var activityTypes []ActivityType
	_, err := this.dbMap.Select(&activityTypes, "select * from activity_types order by title")
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (this *Db) getProjectActivityTypes(projectId int) ([]ActivityType, error) {
	var activityTypes []ActivityType
	_, err := this.dbMap.Select(&activityTypes, "select * from activity_types where id in (select activity_type_id from projects_activity_types where project_id = ?)", projectId)
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (this *Db) deleteProjectActivityTypes(trans *gorp.Transaction, project *Project) error {
	var projectActivityTypes []ProjectActivityTypes
	_, err := trans.Select(&projectActivityTypes, "select * from projects_activity_types where project_id = ?", project.Id)
	if err != nil {
		return err
	}

	itemsToDelete := []ProjectActivityTypes{}
	for _, projectActivityType := range projectActivityTypes {
		deleteItem := true
		for _, activityType := range project.ActivityTypes {
			if activityType.Id == projectActivityType.ActivityTypeId {
				deleteItem = false
				break
			}
		}
		if deleteItem {
			itemsToDelete = append(itemsToDelete, projectActivityType)
		}
	}

	for _, itemToDelete := range itemsToDelete {
		isReferenced, err := this.IsActivityTypeReferenced(itemToDelete.ActivityTypeId)
		if err != nil {
			return err
		}
		if isReferenced {
			return ErrItemInUse
		}
	}

	if len(itemsToDelete) > 0 {
		for _, itemToDelete := range itemsToDelete {
			_, err = trans.Delete(&itemToDelete)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Db) SaveActivityType(activityType *ActivityType) error {
	var err error
	if activityType.Id < 0 {
		err = this.dbMap.Insert(activityType)
	} else {
		_, err = this.dbMap.Update(activityType)
	}
	return err
}

func (this *Db) DeleteActivityType(activityType *ActivityType) error {
	_, err := this.dbMap.Delete(activityType)
	if err != nil {
		return err
	}

	return nil
}

func (this *Db) IsActivityTypeReferenced(id int) (bool, error) {
	exists, err := this.dbMap.SelectInt("select exists(select id from activities where activity_type_id = ?)", id)
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
