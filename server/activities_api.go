package server

import (
	"net/http"
	"strconv"
	"time"
)

type ActivitiesApi struct {
	db *Db
}

func NewActivitiesApi(db *Db) *ActivitiesApi {
	return &ActivitiesApi{db}
}

func (this *ActivitiesApi) GetByDayHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *CtrlHandlerError) {
	dayString := getRouteVar(r, "day")
	day, err := time.Parse("2006-01-02", dayString)
	if err != nil {
		return nil, &CtrlHandlerError{err, "invalid parameter: day", http.StatusBadRequest}
	}

	activities, err := this.getByDay(user.Id, day)
	if err != nil {
		return nil, &CtrlHandlerError{err, "couldn't retrieve activities", http.StatusInternalServerError}
	}
	return activities, nil
}

func (this *ActivitiesApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *CtrlHandlerError) {
	var activity Activity
	unmarshalJson(r.Body, &activity)
	err := this.save(&activity)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultBool(true)
}

func (this *ActivitiesApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *CtrlHandlerError) {
	idString := getRouteVar(r, "id")
	id, err := strconv.ParseInt(idString, 0, 32)
	if err != nil {
		return nil, &CtrlHandlerError{err, "invalid parameter: id", http.StatusBadRequest}
	}

	err = this.delete(int(id))
	if err != nil {
		return nil, &CtrlHandlerError{err, "couldn't delete activity", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (this *ActivitiesApi) getByDay(userId int, day time.Time) ([]Activity, error) {
	activities, err := this.db.GetActivitiesByDay(userId, day)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (this *ActivitiesApi) save(activity *Activity) error {
	var err error
	var existingActivity *Activity

	if activity.Id == -1 {
		existingActivity, err = this.db.TryGetActivity(activity.UserId, activity.Day, activity.Text)
		if err != nil {
			return err
		}
	}

	if existingActivity != nil {
		existingActivity.Duration += activity.Duration
		return this.db.SaveActivity(existingActivity)
	}

	return this.db.SaveActivity(activity)
}

func (this *ActivitiesApi) delete(id int) error {
	activity, err := this.db.GetActivity(id)
	if err != nil {
		return err
	}

	err = this.db.DeleteActivity(activity)
	if err != nil {
		return err
	}

	return nil
}
