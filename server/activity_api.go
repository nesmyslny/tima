package server

import (
	"net/http"
	"time"
)

type ActivityAPI struct {
	db *DB
}

func NewActivityAPI(db *DB) *ActivityAPI {
	return &ActivityAPI{db}
}

func (activityAPI *ActivityAPI) GetByDayHandler(context *HandlerContext) (interface{}, *HandlerError) {
	day, err := context.GetRouteVarTime("day", dateLayout)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	activities, err := activityAPI.getByDay(context.User.ID, day)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve activities", http.StatusInternalServerError}
	}
	return activities, nil
}

func (activityAPI *ActivityAPI) SaveHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var activity Activity
	err := context.GetReqBodyJSON(&activity)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = activityAPI.save(&activity)
	if err != nil {
		return nil, &HandlerError{err, "couldn't save activity", http.StatusInternalServerError}
	}
	return jsonResultBool(true)
}

func (activityAPI *ActivityAPI) DeleteHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = activityAPI.delete(id)
	if err != nil {
		return nil, &HandlerError{err, "couldn't delete activity", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (activityAPI *ActivityAPI) getByDay(userID int, day time.Time) ([]ActivityView, error) {
	activities, err := activityAPI.db.GetActivitiesByDay(userID, day)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (activityAPI *ActivityAPI) save(activity *Activity) error {
	var err error
	var existingActivity *Activity

	if activity.ID == -1 {
		existingActivity, err = activityAPI.db.TryGetActivity(activity.Day, activity.UserID, activity.ProjectID, activity.ActivityTypeID)
		if err != nil {
			return err
		}
	}

	if existingActivity != nil {
		existingActivity.Duration += activity.Duration
		return activityAPI.db.SaveActivity(existingActivity)
	}

	return activityAPI.db.SaveActivity(activity)
}

func (activityAPI *ActivityAPI) delete(id int) error {
	activity, err := activityAPI.db.GetActivity(id)
	if err != nil {
		return err
	}

	err = activityAPI.db.DeleteActivity(activity)
	if err != nil {
		return err
	}

	return nil
}
