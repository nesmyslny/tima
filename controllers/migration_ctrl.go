package controllers

import (
	"github.com/nesmyslny/tima/services"
	"net/http"
)

type MigrationController struct {
	migrationService *services.MigrationService
}

func NewMigrationController(migrationService *services.MigrationService) *MigrationController {
	return &MigrationController{migrationService}
}

func (this *MigrationController) Upgrade(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	err := this.migrationService.Run()
	if err != nil {
		// todo: logging
		// in this case, the internal error is directly exposed to the user.
		// upgrading is an admin task and the internal error is needed to resolve issues.
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}
