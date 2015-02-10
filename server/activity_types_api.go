package server

import "net/http"

type ActivityTypesApi struct {
	db *Db
}

func NewActivityTypesApi(db *Db) *ActivityTypesApi {
	return &ActivityTypesApi{db}
}

func (this *ActivityTypesApi) GetHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	activityType, err := this.get(id)
	if err != nil {
		return nil, &HandlerError{err, "unknown id", http.StatusBadRequest}
	}
	return activityType, nil
}

func (this *ActivityTypesApi) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	activityTypes, err := this.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve activity types", http.StatusInternalServerError}
	}
	return activityTypes, nil
}

func (this *ActivityTypesApi) GetActivityViewListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	list, err := this.getProjectActivityTypesView()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects/activities", http.StatusInternalServerError}
	}
	return list, nil
}

func (this *ActivityTypesApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var activityType ActivityType
	err := unmarshalJson(r.Body, &activityType)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.save(&activityType)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultInt(activityType.Id)
}

func (this *ActivityTypesApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.delete(id)
	if err != nil {
		if err == ErrItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete a activity type that is already in use.", http.StatusBadRequest}
		} else {
			return nil, &HandlerError{err, "Error: Activity type could not deleted.", http.StatusInternalServerError}
		}
	}

	return jsonResultBool(true)
}

func (this *ActivityTypesApi) get(id int) (*ActivityType, error) {
	activityType, err := this.db.GetActivityType(id)
	if err != nil {
		return nil, err
	}
	return activityType, nil
}

func (this *ActivityTypesApi) getList() ([]ActivityType, error) {
	activityTypes, err := this.db.GetActivityTypes()
	if err != nil {
		return nil, err
	}
	return activityTypes, nil
}

func (this *ActivityTypesApi) save(activityType *ActivityType) error {
	return this.db.SaveActivityType(activityType)
}

func (this *ActivityTypesApi) delete(id int) error {
	isReferenced, err := this.db.IsActivityTypeReferenced(id)
	if err != nil {
		return err
	} else if isReferenced {
		return ErrItemInUse
	}

	activityType, err := this.db.GetActivityType(id)
	if err != nil {
		return err
	}

	err = this.db.DeleteActivityType(activityType)
	if err != nil {
		return err
	}

	return nil
}

func (this *ActivityTypesApi) getProjectActivityTypesView() ([]ProjectActivityTypesView, error) {
	return this.db.GetProjectActivityTypesView()
}
