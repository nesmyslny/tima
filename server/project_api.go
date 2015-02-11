package server

import "net/http"

type ProjectApi struct {
	db *Db
}

func NewProjectApi(db *Db) *ProjectApi {
	return &ProjectApi{db}
}

func (this *ProjectApi) GetHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	project, err := this.get(id)
	if err != nil {
		return nil, &HandlerError{err, "unknown id", http.StatusBadRequest}
	}
	return project, nil
}

func (this *ProjectApi) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	projects, err := this.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}
	return projects, nil
}

func (this *ProjectApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var project Project
	err := unmarshalJson(r.Body, &project)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.save(&project)
	if err != nil {
		if err == ErrItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete activity types that are already in use.", http.StatusBadRequest}
		} else {
			return nil, &HandlerError{err, "Error: Project could not be saved.", http.StatusInternalServerError}
		}
	}
	return jsonResultInt(project.Id)
}

func (this *ProjectApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.delete(id)
	if err != nil {
		if err == ErrItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete a project that is already in use.", http.StatusBadRequest}
		} else {
			return nil, &HandlerError{err, "Error: Project could not deleted.", http.StatusInternalServerError}
		}
	}

	return jsonResultBool(true)
}

func (this *ProjectApi) get(id int) (*Project, error) {
	project, err := this.db.GetProject(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (this *ProjectApi) getList() ([]Project, error) {
	projects, err := this.db.GetProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (this *ProjectApi) save(project *Project) error {
	return this.db.SaveProject(project)
}

func (this *ProjectApi) delete(id int) error {
	isReferenced, err := this.db.IsProjectReferenced(id)
	if err != nil {
		return err
	} else if isReferenced {
		return ErrItemInUse
	}

	project, err := this.db.GetProject(id)
	if err != nil {
		return err
	}

	err = this.db.DeleteProject(project)
	if err != nil {
		return err
	}

	return nil
}
