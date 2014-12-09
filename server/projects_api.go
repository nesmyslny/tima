package server

import "net/http"

type ProjectsApi struct {
	db *Db
}

func NewProjectsApi(db *Db) *ProjectsApi {
	return &ProjectsApi{db}
}

func (this *ProjectsApi) GetHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
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

func (this *ProjectsApi) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	projects, err := this.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}
	return projects, nil
}

func (this *ProjectsApi) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var project Project
	err := unmarshalJson(r.Body, &project)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.save(&project)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}
	return jsonResultInt(project.Id)
}

func (this *ProjectsApi) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = this.delete(id)
	if err != nil {
		return nil, &HandlerError{err, "couldn't delete project", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (this *ProjectsApi) get(id int) (*Project, error) {
	project, err := this.db.GetProject(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (this *ProjectsApi) getList() ([]Project, error) {
	projects, err := this.db.GetProjects()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (this *ProjectsApi) save(project *Project) error {
	return this.db.SaveProject(project)
}

func (this *ProjectsApi) delete(id int) error {
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
