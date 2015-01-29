package server

import (
	"net/http"
	"time"
)

type ActivitiesApi struct {
	db *Db
}

func NewActivitiesApi(db *Db) *ActivitiesApi {
	return &ActivitiesApi{db}
}

func (this *ActivitiesApi) GetByDayHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
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

func (this *ActivitiesApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
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

func (this *ActivitiesApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
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

func (this *ActivitiesApi) getByDay(userId int, day time.Time) ([]Activity, error) {
	activities, err := this.db.GetActivitiesByDay(userId, day)
	if err != nil {
		return nil, err
	}
	err = this.setProjectTitle(activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (this *ActivitiesApi) setProjectTitle(activities []Activity) error {
	for i := 0; i < len(activities); i++ {
		projectId := activities[i].ProjectId
		project, err := this.db.GetProject(projectId)
		if err != nil {
			return err
		}
		activities[i].ProjectTitle = project.Title
	}
	return nil
}

func (this *ActivitiesApi) save(activity *Activity) error {
	var err error
	var existingActivity *Activity

	if activity.Id == -1 {
		existingActivity, err = this.db.TryGetActivity(activity.Day, activity.UserId, activity.ProjectId)
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
