package services

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rubenv/sql-migrate"
)

type Migration struct{}

func (this *Migration) Run() error {
	const dialect = "mysql"
	migrate.SetTable("migrations")
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	// todo: configuration of datasource string
	db, err := sql.Open("mysql", "root:pwd@tcp(localhost:3307)/gnomon?parseTime=true")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = migrate.Exec(db, dialect, migrations, migrate.Up)
	if err != nil {
		// todo: logging
		// if(!) any migration were applied, try to roll back:
		// migrate.ExecMax(nil, dialect, migrations, migrate.Down, applied)
		return err
	}

	return nil
}
