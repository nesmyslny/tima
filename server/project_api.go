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

	addedActivityTypes, removedActivityTypes, err := projectAPI.getChangedActivityTypes(project)
	if err != nil {
		return err
	}
	addedUsers, removedUsers, err := projectAPI.getChangedUsers(project)
	if err != nil {
		return err
	}
	return projectAPI.db.SaveProject(project, addedActivityTypes, removedActivityTypes, addedUsers, removedUsers)
}

func (projectAPI *ProjectAPI) getChangedActivityTypes(project *Project) ([]ProjectActivityType, []ProjectActivityType, error) {
	savedProjectActivityTypeIDs, err := projectAPI.db.GetProjectActivityTypeIDs(project.ID)
	if err != nil {
		return nil, nil, err
	}

	projectActivityTypeIDs := project.getActivityTypeIDs()
	addedProjectActivityTypes := projectAPI.getAddedActivityTypes(project.ID, projectActivityTypeIDs, savedProjectActivityTypeIDs)
	removedActivityTypes, err := projectAPI.getRemovedActivityTypes(project.ID, projectActivityTypeIDs, savedProjectActivityTypeIDs)
	if err != nil {
		return nil, nil, err
	}

	return addedProjectActivityTypes, removedActivityTypes, nil
}

func (projectAPI *ProjectAPI) getAddedActivityTypes(projectID int, projectActivityTypeIDs, savedProjectActivityTypeIDs []int) []ProjectActivityType {
	addedIDs := diffInt(projectActivityTypeIDs, savedProjectActivityTypeIDs)
	addedProjectActivityTypes, _ := createProjectActivityTypes(projectID, addedIDs, nil)
	return addedProjectActivityTypes
}

func (projectAPI *ProjectAPI) getRemovedActivityTypes(projectID int, projectActivityTypeIDs, savedProjectActivityTypeIDs []int) ([]ProjectActivityType, error) {
	removedIDs := diffInt(savedProjectActivityTypeIDs, projectActivityTypeIDs)
	return createProjectActivityTypes(projectID, removedIDs, projectAPI.createRemovedActivityTypesCheck)
}

func (projectAPI *ProjectAPI) createRemovedActivityTypesCheck(projectID, activityTypeID int) error {
	isReferenced, err := projectAPI.db.IsActivityTypeReferenced(activityTypeID, &projectID)
	if err != nil {
		return err
	}
	if isReferenced {
		return errItemInUse
	}
	return nil
}

func (projectAPI *ProjectAPI) getChangedUsers(project *Project) ([]ProjectUser, []ProjectUser, error) {
	savedUserIDs, err := projectAPI.db.GetProjectUserIDs(project.ID)
	if err != nil {
		return nil, nil, err
	}

	userIDs := project.GetUserIDs()
	addedProjectUsers := projectAPI.getAddedUsers(project.ID, userIDs, savedUserIDs)
	removedProjectUsers := projectAPI.getRemovedUsers(project.ID, userIDs, savedUserIDs)

	return addedProjectUsers, removedProjectUsers, nil
}

func (projectAPI *ProjectAPI) getAddedUsers(projectID int, userIDs, savedUserIDs []int) []ProjectUser {
	addedIDs := diffInt(userIDs, savedUserIDs)
	return createProjectUsers(projectID, addedIDs)
}

func (projectAPI *ProjectAPI) getRemovedUsers(projectID int, userIDs, savedUserIDs []int) []ProjectUser {
	removedIDs := diffInt(savedUserIDs, userIDs)
	return createProjectUsers(projectID, removedIDs)
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

func createProjectUsers(projectID int, userIDs []int) []ProjectUser {
	var projectUsers []ProjectUser
	for _, userID := range userIDs {
		projectUsers = append(projectUsers, ProjectUser{projectID, userID})
	}
	return projectUsers
}

func createProjectActivityTypes(projectID int, activityTypeIDs []int, checkCallback func(int, int) error) ([]ProjectActivityType, error) {
	var projectActivityTypes []ProjectActivityType
	for _, activityTypeID := range activityTypeIDs {
		if checkCallback != nil {
			if err := checkCallback(projectID, activityTypeID); err != nil {
				return nil, err
			}
		}
		projectActivityTypes = append(projectActivityTypes, ProjectActivityType{projectID, activityTypeID})
	}
	return projectActivityTypes, nil
}
