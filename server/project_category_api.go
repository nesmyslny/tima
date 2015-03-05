package server

import "net/http"

type ProjectCategoryAPI struct {
	db *DB
}

func NewProjectCategoryAPI(db *DB) *ProjectCategoryAPI {
	return &ProjectCategoryAPI{db}
}

func (projectCategoryAPI *ProjectCategoryAPI) GetTreeHandler(context *HandlerContext) (interface{}, *HandlerError) {
	projectCategories, err := projectCategoryAPI.getTree(nil)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve project categories", http.StatusInternalServerError}
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) GetListHandler(context *HandlerContext) (interface{}, *HandlerError) {
	projectCategories, err := projectCategoryAPI.getList(nil)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve project categories", http.StatusInternalServerError}
	}
	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) SaveHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var projectCategory ProjectCategory
	err := context.GetReqBodyJSON(&projectCategory)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectCategoryAPI.save(&projectCategory)
	if err != nil {
		if err == errIDNotUnique {
			return nil, &HandlerError{err, "Error: Reference ID is already in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Project category could not be saved.", http.StatusInternalServerError}
	}
	return jsonResultInt(projectCategory.ID)
}

func (projectCategoryAPI *ProjectCategoryAPI) DeleteHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
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

func (projectCategoryAPI *ProjectCategoryAPI) getTree(parent *ProjectCategory) ([]ProjectCategory, error) {
	projectCategories, err := projectCategoryAPI.db.GetProjectCategories(parent)
	if err != nil {
		return nil, err
	}

	for i := range projectCategories {
		projectCategories[i].ProjectCategories, err = projectCategoryAPI.getTree(&projectCategories[i])
		if err != nil {
			return nil, err
		}
	}

	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) getList(parent *ProjectCategory) ([]ProjectCategory, error) {
	projectCategories, err := projectCategoryAPI.db.GetProjectCategories(parent)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(projectCategories); i++ {
		children, err := projectCategoryAPI.getList(&projectCategories[i])
		if err != nil {
			return nil, err
		}

		// inserting children into slice after parent
		slicingIndex := i + 1
		projectCategories = append(projectCategories[:slicingIndex], append(children, projectCategories[slicingIndex:]...)...)
		i += len(children)
	}

	return projectCategories, nil
}

func (projectCategoryAPI *ProjectCategoryAPI) save(projectCategory *ProjectCategory) error {
	return projectCategoryAPI.db.SaveProjectCategory(projectCategory)
}

func (projectCategoryAPI *ProjectCategoryAPI) delete(id int) error {
	projectCatetory, err := projectCategoryAPI.db.GetProjectCategory(id)
	if err != nil {
		return err
	}

	isReferenced, err := projectCategoryAPI.db.IsProjectCategoryReferenced(projectCatetory)
	if err != nil {
		return err
	} else if isReferenced {
		return errItemInUse
	}

	err = projectCategoryAPI.db.DeleteProjectCategory(projectCatetory)
	if err != nil {
		return err
	}

	return nil
}
