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
		if err == errOptimisticLocking {
			return nil, &HandlerError{err, "Error: The activity was changed/deleted by another user.", http.StatusInternalServerError}
		}
		return nil, &HandlerError{err, "couldn't save activity", http.StatusInternalServerError}
	}
	return activity, nil
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

	return &SingleValue{true}, nil
}

func (activityAPI *ActivityAPI) getByDay(userID int, day time.Time) ([]ActivityView, error) {
	activities, err := activityAPI.db.GetActivitiesByDay(userID, day)
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (activityAPI *ActivityAPI) save(activity *Activity) error {
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
