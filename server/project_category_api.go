package server

import "net/http"

type ProjectCategoryAPI struct {
	db *DB
}

func NewProjectCategoryAPI(db *DB) *ProjectCategoryAPI {
	return &ProjectCategoryAPI{db}
}

func (projectCategoryAPI *ProjectCategoryAPI) GetTreeHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	projectCategories, err := projectCategoryAPI.getTree()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve project categories", http.StatusInternalServerError}
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) GetListHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	projectCategories, err := projectCategoryAPI.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve project categories", http.StatusInternalServerError}
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) SaveHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	var projectCategory ProjectCategory
	err := unmarshalJSON(r.Body, &projectCategory)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectCategoryAPI.save(&projectCategory)
	if err != nil {
		return nil, &HandlerError{err, "Error: Project category could not be saved.", http.StatusInternalServerError}
	}
	return jsonResultInt(projectCategory.ID)
}

func (projectCategoryAPI *ProjectCategoryAPI) DeleteHandler(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError) {
	id, err := getRouteVarInt(r, "id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectCategoryAPI.delete(id)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: This project category, or one of its descendants, is in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Project category could not deleted.", http.StatusInternalServerError}
	}

	return jsonResultBool(true)
}

func (projectCategoryAPI *ProjectCategoryAPI) getTree() ([]ProjectCategory, error) {
	projectCategories, err := projectCategoryAPI.db.GetProjectCategoryTree(nil)
	if err != nil {
		return nil, err
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) getList() ([]ProjectCategory, error) {
	projectCategories, err := projectCategoryAPI.db.GetProjectCategoryList(nil)
	if err != nil {
		return nil, err
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) save(projectCategory *ProjectCategory) error {
	return projectCategoryAPI.db.SaveProjectCategory(projectCategory)
}

func (projectCategoryAPI *ProjectCategoryAPI) delete(id int) error {
	isReferenced, err := projectCategoryAPI.db.IsProjectCategoryReferenced(id)
	if err != nil {
		return err
	} else if isReferenced {
		return errItemInUse
	}

	projectCatetory, err := projectCategoryAPI.db.GetProjectCategory(id)
	if err != nil {
		return err
	}

	err = projectCategoryAPI.db.DeleteProjectCategory(projectCatetory)
	if err != nil {
		return err
	}

	return nil
}
