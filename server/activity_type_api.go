package server

import "net/http"

type ActivityTypeAPI struct {
	db *DB
}

func NewActivityTypeAPI(db *DB) *ActivityTypeAPI {
	return &ActivityTypeAPI{db}
}

func (activityTypeAPI *ActivityTypeAPI) GetHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	activityType, err := activityTypeAPI.get(id)
	if err != nil {
		return nil, &HandlerError{err, "unknown id", http.StatusBadRequest}
	}
	return activityType, nil
}

func (activityTypeAPI *ActivityTypeAPI) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	activityTypes, err := activityTypeAPI.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve activity types", http.StatusInternalServerError}
	}
	return activityTypes, nil
}

func (activityTypeAPI *ActivityTypeAPI) GetActivityViewListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	list, err := activityTypeAPI.getProjectActivityTypeViewList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects/activities", http.StatusInternalServerError}
	}
	return list, nil
}

func (activityTypeAPI *ActivityTypeAPI) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var activityType ActivityType
	err := unmarshalJSON(r.Body, &activityType)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = activityTypeAPI.save(&activityType)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultInt(activityType.ID)
}

func (activityTypeAPI *ActivityTypeAPI) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = activityTypeAPI.delete(id)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete a activity type that is already in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Activity type could not deleted.", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (activityTypeAPI *ActivityTypeAPI) get(id int) (*ActivityType, error) {
	activityType, err := activityTypeAPI.db.GetActivityType(id)
	if err != nil {
		return nil, err
	}
	return activityType, nil
}

func (activityTypeAPI *ActivityTypeAPI) getList() ([]ActivityType, error) {
	activityTypes, err := activityTypeAPI.db.GetActivityTypes()
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (activityTypeAPI *ActivityTypeAPI) save(activityType *ActivityType) error {
	return activityTypeAPI.db.SaveActivityType(activityType)
}

func (activityTypeAPI *ActivityTypeAPI) delete(id int) error {
	isReferenced, err := activityTypeAPI.db.IsActivityTypeReferenced(id, nil)
	if err != nil {
		return err
	} else if isReferenced {
		return errItemInUse
	}

	activityType, err := activityTypeAPI.db.GetActivityType(id)
	if err != nil {
		return err
	}

	err = activityTypeAPI.db.DeleteActivityType(activityType)
	if err != nil {
		return err
	}

	return nil
}

func (activityTypeAPI *ActivityTypeAPI) getProjectActivityTypeViewList() ([]ProjectActivityTypeView, error) {
	return activityTypeAPI.db.GetProjectActivityTypeViewList()
}
