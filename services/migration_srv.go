package services

import (
	"gnomon/dbaccess"
	"gnomon/models"
)

type MigrationService struct {
	db *DbAccess.Db
}

func NewMigrationService(db *DbAccess.Db) *MigrationService {
	return &MigrationService{db}
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
		user := &models.User{
			Id:           -1,
			Username:     "admin",
			PasswordHash: "pwd",
		}
		this.db.SaveUser(user)
	}

	return nil
}
