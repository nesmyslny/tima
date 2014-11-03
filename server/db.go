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

func (this *Db) TryGetActivity(userId int, day time.Time, text string) (*Activity, error) {
	var activity *Activity
	err := this.dbMap.SelectOne(&activity,
		"select * from activities where user_id = ? and day = ? and text = ?",
		userId, day.Format(dateLayout), text)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return activity, nil
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
	return obj.(*Project), nil
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
	var err error
	if project.Id < 0 {
		err = this.dbMap.Insert(project)
	} else {
		_, err = this.dbMap.Update(project)
	}
	return err
}
