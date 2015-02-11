package server

import (
	"net/http"
	"time"
)

type ActivityApi struct {
	db *Db
}

func NewActivityApi(db *Db) *ActivityApi {
	return &ActivityApi{db}
}

func (this *ActivityApi) GetByDayHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	day, err := getRouteVarTime(r, "day", "2006-01-02")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	activities, err := this.getByDay(user.Id, day)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve activities", http.StatusInternalServerError}
	}
	return activities, nil
}

func (this *ActivityApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var activity Activity
	err := unmarshalJson(r.Body, &activity)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.save(&activity)
	if err != nil {
		return nil, &HandlerError{err, "couldn't save activity", http.StatusInternalServerError}
	}
	return jsonResultBool(true)
}

func (this *ActivityApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.delete(id)
	if err != nil {
		return nil, &HandlerError{err, "couldn't delete activity", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (this *ActivityApi) getByDay(userId int, day time.Time) ([]ActivityView, error) {
	activities, err := this.db.GetActivitiesByDay(userId, day)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (this *ActivityApi) save(activity *Activity) error {
	var err error
	var existingActivity *Activity

	if activity.Id == -1 {
		existingActivity, err = this.db.TryGetActivity(activity.Day, activity.UserId, activity.ProjectId, activity.ActivityTypeId)
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

func (this *ActivityApi) delete(id int) error {
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
