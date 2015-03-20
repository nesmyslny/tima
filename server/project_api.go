package server

import "net/http"

type ProjectAPI struct {
	db *DB
}

func NewProjectAPI(db *DB) *ProjectAPI {
	return &ProjectAPI{db}
}

func (projectAPI *ProjectAPI) authorizeGetSave(projectID int, user *User) (bool, error) {
	project, err := projectAPI.db.GetProject(projectID)
	if err != nil {
		return false, err
	}
	return *user.Role >= RoleManager || *project.ResponsibleUserID == user.ID || *project.ManagerUserID == user.ID, nil
}

func (projectAPI *ProjectAPI) AuthorizeGet(context *HandlerContext) (bool, error) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return false, err
	}
	return projectAPI.authorizeGetSave(id, context.User)
}

func (projectAPI *ProjectAPI) AuthorizeSave(context *HandlerContext) (bool, error) {
	var project Project
	err := context.GetReqBodyJSON(&project)
	if err != nil {
		return false, err
	}
	return projectAPI.authorizeGetSave(project.ID, context.User)
}

func (projectAPI *ProjectAPI) GetHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	project, err := projectAPI.get(id)
	if err != nil {
		return nil, &HandlerError{err, "unknown id", http.StatusBadRequest}
	}
	return project, nil
}

func (projectAPI *ProjectAPI) GetListHandler(context *HandlerContext) (interface{}, *HandlerError) {
	projects, err := projectAPI.getList()
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) GetListUserHandler(context *HandlerContext) (interface{}, *HandlerError) {
	projects, err := projectAPI.getListUser(context.User.ID)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) SaveHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var project Project
	err := context.GetReqBodyJSON(&project)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = projectAPI.save(&project, context.User)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: It is not possible to delete activity types that are already in use.", http.StatusBadRequest}
		} else if err == errIDNotUnique {
			return nil, &HandlerError{err, "Error: Reference ID is already in use.", http.StatusBadRequest}
		} else if err == errOptimisticLocking {
			return nil, &HandlerError{err, "Error: The project was changed/deleted by another user.", http.StatusInternalServerError}
		}
		return nil, &HandlerError{err, "Error: Project could not be saved.", http.StatusInternalServerError}
	}
	return project, nil
}

func (projectAPI *ProjectAPI) DeleteHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
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

	return &SingleValue{true}, nil
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

func (projectAPI *ProjectAPI) getListUser(userID int) ([]Project, error) {
	projects, err := projectAPI.db.GetProjectsOfUser(userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) save(project *Project, user *User) error {
	// creating projects is only allowd for admins and managers
	// updating projects is partially allowed for users, who are responsible or manager of the project
	if project.ID < 0 && *user.Role < RoleManager {
		return errForbidden
	} else if project.ID >= 0 && *user.Role < RoleManager {
		projectOrig, err := projectAPI.db.GetProject(project.ID)
		if *projectOrig.ResponsibleUserID == user.ID || *projectOrig.ManagerUserID == user.ID {
			if err != nil {
				return err
			}

			project.RefID = projectOrig.RefID
			project.RefIDComplete = projectOrig.RefIDComplete
			project.Title = projectOrig.Title
			project.ProjectCategoryID = projectOrig.ProjectCategoryID
			project.ResponsibleUserID = projectOrig.ResponsibleUserID

			// when the user is only project manager, changing manager is not allowed
			if *projectOrig.ResponsibleUserID != user.ID && *projectOrig.ManagerUserID == user.ID {
				project.ManagerUserID = projectOrig.ManagerUserID
			}
		} else {
			return errForbidden
		}
	}

	addedActivityTypes, removedActivityItems, err := projectAPI.getChangedActivityTypes(project)
	if err != nil {
		return err
	}
	addedUsers, removedUsers, err := projectAPI.getChangedUsers(project)
	if err != nil {
		return err
	}
	return projectAPI.db.SaveProject(project, addedActivityTypes, removedActivityItems, addedUsers, removedUsers)
}

func (projectAPI *ProjectAPI) getChangedActivityTypes(project *Project) ([]ProjectActivityType, []ProjectActivityType, error) {
	projectActivityTypes, err := projectAPI.db.GetProjectActivityTypes(project.ID)
	if err != nil {
		return nil, nil, err
	}

	addedItems, err := projectAPI.getAddedActivityTypes(project, projectActivityTypes)
	if err != nil {
		return nil, nil, err
	}

	removedItems, err := projectAPI.getRemovedActivityTypes(project, projectActivityTypes)
	if err != nil {
		return nil, nil, err
	}

	return addedItems, removedItems, nil
}

func (projectAPI *ProjectAPI) getAddedActivityTypes(project *Project, projectActivityTypes []ProjectActivityType) ([]ProjectActivityType, error) {
	var addedItems []ProjectActivityType
	for _, activityType := range project.ActivityTypes {
		added := true
		for _, projectActivityType := range projectActivityTypes {
			if projectActivityType.ActivityTypeID == activityType.ID {
				added = false
				break
			}
		}

		if added {
			addedItems = append(addedItems, ProjectActivityType{project.ID, activityType.ID})
		}
	}

	return addedItems, nil
}

func (projectAPI *ProjectAPI) getRemovedActivityTypes(project *Project, projectActivityTypes []ProjectActivityType) ([]ProjectActivityType, error) {
	var removedItems []ProjectActivityType
	for _, projectActivityType := range projectActivityTypes {
		deleted := true
		for _, activityType := range project.ActivityTypes {
			if activityType.ID == projectActivityType.ActivityTypeID {
				deleted = false
				break
			}
		}

		if deleted {
			isReferenced, err := projectAPI.db.IsActivityTypeReferenced(projectActivityType.ActivityTypeID, &project.ID)
			if err != nil {
				return nil, err
			}
			if isReferenced {
				return nil, errItemInUse
			}

			removedItems = append(removedItems, projectActivityType)
		}
	}

	return removedItems, nil
}

func (projectAPI *ProjectAPI) getChangedUsers(project *Project) ([]ProjectUser, []ProjectUser, error) {
	projectUsers, err := projectAPI.db.GetProjectUsers(project.ID)
	if err != nil {
		return nil, nil, err
	}

	addedItems, err := projectAPI.getAddedUsers(project, projectUsers)
	if err != nil {
		return nil, nil, err
	}

	removedItems, err := projectAPI.getRemovedUsers(project, projectUsers)
	if err != nil {
		return nil, nil, err
	}

	return addedItems, removedItems, nil
}

func (projectAPI *ProjectAPI) getAddedUsers(project *Project, projectUsers []ProjectUser) ([]ProjectUser, error) {
	var addedItems []ProjectUser
	for _, user := range project.Users {
		added := true
		for _, projectUser := range projectUsers {
			if projectUser.UserID == user.ID {
				added = false
				break
			}
		}

		if added {
			addedItems = append(addedItems, ProjectUser{project.ID, user.ID})
		}
	}

	return addedItems, nil
}

func (projectAPI *ProjectAPI) getRemovedUsers(project *Project, projectUsers []ProjectUser) ([]ProjectUser, error) {
	var removedItems []ProjectUser
	for _, projectUser := range projectUsers {
		deleted := true
		for _, user := range project.Users {
			if user.ID == projectUser.UserID {
				deleted = false
				break
			}
		}

		if deleted {
			removedItems = append(removedItems, projectUser)
		}
	}

	return removedItems, nil
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
