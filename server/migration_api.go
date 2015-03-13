package server

import (
	"net/http"
)

type MigrationAPI struct {
	db      *DB
	userAPI *UserAPI
}

func NewMigrationAPI(db *DB, userAPI *UserAPI) *MigrationAPI {
	return &MigrationAPI{db, userAPI}
}

func (migrationAPI *MigrationAPI) UpgradeHandler(context *HandlerContext) (interface{}, *HandlerError) {
	err := migrationAPI.migrate()
	if err != nil {
		// todo: logging
		// in this case, the internal error is directly exposed to the user.
		// upgrading is an admin task and the internal error is needed to resolve issues.
		return nil, &HandlerError{err, err.Error(), http.StatusInternalServerError}
	}

	return &SingleValue{true}, nil
}

func (migrationAPI *MigrationAPI) migrate() error {
	err := migrationAPI.db.Upgrade(0)
	if err != nil {
		return err
	}

	err = migrationAPI.postMigration()
	if err != nil {
		return err
	}

	return nil
}

func (migrationAPI *MigrationAPI) postMigration() error {
	countUsers, err := migrationAPI.db.GetNumberOfUsers()
	if err != nil {
		return err
	}

	if countUsers == 0 {
		_, err = migrationAPI.userAPI.AddUser("admin", RoleAdmin, 0, "pwd", "", "", "")
	}

	return err
}
