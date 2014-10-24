package services

import (
	"github.com/nesmyslny/tima/dbaccess"
)

type MigrationService struct {
	db          *DbAccess.Db
	userService *UserService
}

func NewMigrationService(db *DbAccess.Db, userService *UserService) *MigrationService {
	return &MigrationService{db, userService}
}

func (this *MigrationService) Run() error {
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

func (this *MigrationService) postMigration() error {
	countUsers, err := this.db.GetNumberOfUsers()
	if err != nil {
		return err
	}

	if countUsers == 0 {
		_, err = this.userService.AddUser("admin", "pwd", "", "", "")
	}

	return err
}
