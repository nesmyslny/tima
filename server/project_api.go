package server

import "net/http"

type ProjectAPI struct {
	db *DB
}

func NewProjectAPI(db *DB) *ProjectAPI {
	return &ProjectAPI{db}
}

func (projectAPI *ProjectAPI) GetHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	project, err := projectAPI.get(id)
	if err != nil {
		return nil, &HandlerError{err, "unknown id", http.StatusBadRequest}
	}
	return project, nil
}

func (projectAPI *ProjectAPI) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	projects, err := projectAPI.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var project Project
	err := unmarshalJSON(r.Body, &project)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectAPI.save(&project)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete activity types that are already in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Project could not be saved.", http.StatusInternalServerError}
	}
	return jsonResultInt(project.ID)
}

func (projectAPI *ProjectAPI) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectAPI.delete(id)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete a project that is already in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Project could not deleted.", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (projectAPI *ProjectAPI) get(id int) (*Project, error) {
	project, err := projectAPI.db.GetProject(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (projectAPI *ProjectAPI) getList() ([]Project, error) {
	projects, err := projectAPI.db.GetProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) save(project *Project) error {
	return projectAPI.db.SaveProject(project)
}

func (projectAPI *ProjectAPI) delete(id int) error {
	isReferenced, err := projectAPI.db.IsProjectReferenced(id)
	if err != nil {
		return err
	} else if isReferenced {
		return errItemInUse
	}

	project, err := projectAPI.db.GetProject(id)
	if err != nil {
		return err
	}

	err = projectAPI.db.DeleteProject(project)
	if err != nil {
		return err
	}

	return nil
}
