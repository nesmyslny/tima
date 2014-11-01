package server

import (
	"net/http"
)

type MigrationApi struct {
	db      *Db
	userApi *UserApi
}

func NewMigrationApi(db *Db, userApi *UserApi) *MigrationApi {
	return &MigrationApi{db, userApi}
}

func (this *MigrationApi) UpgradeHandler(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	err := this.migrate()
	if err != nil {
		// todo: logging
		// in this case, the internal error is directly exposed to the user.
		// upgrading is an admin task and the internal error is needed to resolve issues.
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (this *MigrationApi) migrate() error {
	err := this.db.Upgrade()
	if err != nil {
		return err
	}

	err = this.postMigration()
	if err != nil {
		return err
	}

	return nil
}

func (this *MigrationApi) postMigration() error {
	countUsers, err := this.db.GetNumberOfUsers()
	if err != nil {
		return err
	}

	if countUsers == 0 {
		_, err = this.userApi.AddUser("admin", "pwd", "", "", "")
	}

	return err
}
