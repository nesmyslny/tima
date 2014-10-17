package controllers

import (
	"gnomon/services"
	"net/http"
)

type MigrationController struct {
	migrationService *services.MigrationService
}

func NewMigrationController(migrationService *services.MigrationService) *MigrationController {
	return &MigrationController{migrationService}
}

func (this *MigrationController) Upgrade(w http.ResponseWriter, r *http.Request) {
	err := this.migrationService.Run()
	if err != nil {
		// todo: logging
		// in this case, the internal error is directly exposed to the user.
		// upgrading is an admin task and the internal error is needed to resolve issues.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
