package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/nesmyslny/tima/models"
	"github.com/nesmyslny/tima/services"
)

type ActivitiesController struct {
	activitiesService *services.ActivitiesService
}

func NewActivitiesController(activitiesService *services.ActivitiesService) *ActivitiesController {
	return &ActivitiesController{activitiesService}
}

func (this *ActivitiesController) GetActivities(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError) {
	dayString := getRouteVar(r, "day")
	day, err := time.Parse("2006-01-02", dayString)
	if err != nil {
		return nil, &CtrlHandlerError{err, "invalid parameter: day", http.StatusBadRequest}
	}

	activities, err := this.activitiesService.GetActivities(user.Id, day)
	if err != nil {
		return nil, &CtrlHandlerError{err, "couldn't retrieve activities", http.StatusInternalServerError}
	}
	return activities, nil
}

func (this *ActivitiesController) SaveActivity(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError) {
	var activity models.Activity
	unmarshalJson(r.Body, &activity)
	err := this.activitiesService.SaveActivity(&activity)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultBool(true)
}

func (this *ActivitiesController) DeleteActivity(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError) {
	idString := getRouteVar(r, "id")
	id, err := strconv.ParseInt(idString, 0, 32)
	if err != nil {
		return nil, &CtrlHandlerError{err, "invalid parameter: id", http.StatusBadRequest}
	}

	err = this.activitiesService.DeleteActivity(int(id))
	if err != nil {
		return nil, &CtrlHandlerError{err, "couldn't delete activity", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}
