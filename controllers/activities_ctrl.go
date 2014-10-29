package controllers

import (
	"log"
	"net/http"
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

func (this *ActivitiesController) AddActivity(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError) {
	var activity models.Activity
	unmarshalJson(r.Body, &activity)
	log.Print(activity.Day)
	activity.UserId = user.Id
	err := this.activitiesService.AddActivity(&activity)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultBool(true)
}
