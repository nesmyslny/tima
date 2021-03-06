package server

import "net/http"

type ProjectAPI struct {
	db *DB
}

func NewProjectAPI(db *DB) *ProjectAPI {
	return &ProjectAPI{db}
}

func (projectAPI *ProjectAPI) authorizeGetSave(project *Project, user *User) (bool, error) {
	// Authorized are...
	// ...users with role 'Manager'.
	// ...the user that is responsible for the project.
	// ...the user that is the manager for the project.
	if *user.Role >= RoleManager || *project.ResponsibleUserID == user.ID || *project.ManagerUserID == user.ID {
		return true, nil
	}

	// ...users with role 'Department Manager', if the project is in the same department (or a child of the users' department).
	// (the department of a project is defined by the responsible user.)
	deptManAuth, err := projectAPI.isDeptManagerAuthorized(project, user)
	if err != nil {
		return false, err
	}

	if deptManAuth {
		return true, nil
	}

	return false, nil
}

func (projectAPI *ProjectAPI) AuthorizeGet(context *HandlerContext) (bool, error) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return false, err
	}

	project, err := projectAPI.db.GetProject(id)
	if err != nil {
		return false, err
	}

	return projectAPI.authorizeGetSave(project, context.User)
}

func (projectAPI *ProjectAPI) AuthorizeSave(context *HandlerContext) (bool, error) {
	var project Project
	err := context.GetReqBodyJSON(&project)
	if err != nil {
		return false, err
	}
	return projectAPI.authorizeGetSave(&project, context.User)
}

func (projectAPI *ProjectAPI) AuthorizeDelete(context *HandlerContext) (bool, error) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return false, err
	}

	project, err := projectAPI.db.GetProject(id)
	if err != nil {
		return false, err
	}

	if *context.User.Role >= RoleManager {
		return true, nil
	}

	deptManAuth, err := projectAPI.isDeptManagerAuthorized(project, context.User)
	if err != nil {
		return false, err
	}

	if deptManAuth {
		return true, nil
	}

	return false, nil
}

func (projectAPI *ProjectAPI) isDeptManagerAuthorized(project *Project, user *User) (bool, error) {
	// ...users with role 'Department Manager', if the project is in the same department (or a child of the users' department).
	// (the department of a project is defined by the responsible user.)
	if *user.Role == RoleDeptManager {
		isDeptProject, err := projectAPI.isProjectOfDepartment(project, *user.DepartmentID)
		if err != nil {
			return false, err
		}

		if isDeptProject {
			return true, nil
		}
	}

	return false, nil
}

// Checks if the project is of the given department.
// Project don't have a department assign directly. The department of a project is defined by the responsible user.
func (projectAPI *ProjectAPI) isProjectOfDepartment(project *Project, departmentID int) (bool, error) {
	deptIDs, err := projectAPI.db.GetDepartmentIDsDownward(departmentID)
	if err != nil {
		return false, err
	}

	responsibleUser, err := projectAPI.db.GetUser(*project.ResponsibleUserID)
	if err != nil {
		return false, err
	}

	for _, id := range deptIDs {
		if *responsibleUser.DepartmentID == id {
			return true, nil
		}
	}

	return false, nil
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
	var deptID *int
	// department manager are only allowed to edit/view the project of their department (or a descendant department).
	if *context.User.Role == RoleDeptManager {
		deptID = context.User.DepartmentID
	}

	projects, err := projectAPI.getList(deptID)
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

func (projectAPI *ProjectAPI) GetListSelectHandler(context *HandlerContext) (interface{}, *HandlerError) {
	projects, err := projectAPI.getListSelect(context.User.ID, *context.User.DepartmentID)
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

func (projectAPI *ProjectAPI) getList(deptID *int) ([]Project, error) {
	projects, err := projectAPI.db.GetProjects(deptID)
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

func (projectAPI *ProjectAPI) getListSelect(userID int, departmentID int) ([]Project, error) {
	projects, err := projectAPI.db.GetProjectsForSelection(userID, departmentID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (projectAPI *ProjectAPI) save(project *Project, user *User) error {
	err := projectAPI.prepareSave(project, user)
	if err != nil {
		return err
	}

	addedActivityTypes, removedActivityTypes, err := projectAPI.getChangedActivityTypes(project)
	if err != nil {
		return err
	}
	addedUsers, removedUsers, err := projectAPI.getChangedUsers(project)
	if err != nil {
		return err
	}
	addedDepartments, removedDepartments, err := projectAPI.getChangedDepartments(project)
	if err != nil {
		return err
	}
	return projectAPI.db.SaveProject(project, addedActivityTypes, removedActivityTypes, addedDepartments, removedDepartments, addedUsers, removedUsers)
}

func (projectAPI *ProjectAPI) prepareSave(project *Project, user *User) error {
	// creating projects is only allowd for admins, managers and department manager
	// updating projects is partially allowed for users, who are responsible or manager of the project
	if project.ID < 0 && *user.Role < RoleDeptManager {
		return errForbidden
	} else if project.ID >= 0 && *user.Role < RoleManager {
		projectSpecificPrivelege := true

		projectOrig, err := projectAPI.db.GetProject(project.ID)
		if err != nil {
			return err
		}

		// department manager are allowed to save all attributes in projects, that are in their (or a descendant) department
		// if the project is not in their department, the project specific privileges are crucial
		if *user.Role == RoleDeptManager {
			isDeptProject, err := projectAPI.isProjectOfDepartment(projectOrig, *user.DepartmentID)
			if err != nil {
				return err
			}

			if isDeptProject {
				projectOfDept, err := projectAPI.isProjectOfDepartment(project, *user.DepartmentID)
				if err != nil {
					return err
				}
				if !projectOfDept {
					return errForbidden
				}

				projectSpecificPrivelege = false
			}
		}

		if projectSpecificPrivelege {
			if *projectOrig.ResponsibleUserID == user.ID || *projectOrig.ManagerUserID == user.ID {
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
	}

	return nil
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

func (projectAPI *ProjectAPI) getChangedDepartments(project *Project) ([]ProjectDepartment, []ProjectDepartment, error) {
	savedDepartmentIDs, err := projectAPI.db.GetProjectDepartmentIDs(project.ID)
	if err != nil {
		return nil, nil, err
	}

	departmentIDs := project.GetDepartmentIDs()
	addedProjectDepartments := projectAPI.getAddedDepartments(project.ID, departmentIDs, savedDepartmentIDs)
	removedProjectDepartments := projectAPI.getRemovedDepartments(project.ID, departmentIDs, savedDepartmentIDs)

	return addedProjectDepartments, removedProjectDepartments, nil
}

func (projectAPI *ProjectAPI) getAddedDepartments(projectID int, departmentIDs, savedDepartmentIDs []int) []ProjectDepartment {
	addedIDs := diffInt(departmentIDs, savedDepartmentIDs)
	return createProjectDepartments(projectID, addedIDs)
}

func (projectAPI *ProjectAPI) getRemovedDepartments(projectID int, departmentIDs, savedDepartmentIDs []int) []ProjectDepartment {
	removedIDs := diffInt(savedDepartmentIDs, departmentIDs)
	return createProjectDepartments(projectID, removedIDs)
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

func createProjectDepartments(projectID int, departmentIDs []int) []ProjectDepartment {
	var projectDepartments []ProjectDepartment
	for _, departmentID := range departmentIDs {
		projectDepartments = append(projectDepartments, ProjectDepartment{projectID, departmentID})
	}
	return projectDepartments
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
