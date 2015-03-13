package server

import "net/http"

type DepartmentAPI struct {
	db *DB
}

func NewDepartmentAPI(db *DB) *DepartmentAPI {
	return &DepartmentAPI{db}
}

func (departmentAPI *DepartmentAPI) GetTreeHandler(context *HandlerContext) (interface{}, *HandlerError) {
	departments, err := departmentAPI.getTree(nil)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve departments", http.StatusInternalServerError}
	}
	return departments, nil
}

func (departmentAPI *DepartmentAPI) GetListHandler(context *HandlerContext) (interface{}, *HandlerError) {
	departments, err := departmentAPI.getList(nil)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve departments", http.StatusInternalServerError}
	}
	return departments, nil
}

func (departmentAPI *DepartmentAPI) SaveHandler(context *HandlerContext) (interface{}, *HandlerError) {
	var department Department
	err := context.GetReqBodyJSON(&department)
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = departmentAPI.save(&department)
	if err != nil {
		if err == errOptimisticLocking {
			return nil, &HandlerError{err, "Error: Department was changed/deleted by another user.", http.StatusInternalServerError}
		}
		return nil, &HandlerError{err, "Error: Department could not be saved.", http.StatusInternalServerError}
	}
	return department, nil
}

func (departmentAPI *DepartmentAPI) DeleteHandler(context *HandlerContext) (interface{}, *HandlerError) {
	id, err := context.GetRouteVarInt("id")
	if err != nil {
		return nil, &HandlerError{err, err.Error(), http.StatusBadRequest}
	}

	err = departmentAPI.delete(id)
	if err != nil {
		if err == errItemInUse {
			return nil, &HandlerError{err, "Error: This department, or one of its descendants, is in use.", http.StatusBadRequest}
		}
		return nil, &HandlerError{err, "Error: Department could not deleted.", http.StatusInternalServerError}
	}

	return &SingleValue{true}, nil
}

func (departmentAPI *DepartmentAPI) getTree(parent *Department) ([]Department, error) {
	departments, err := departmentAPI.db.GetDepartments(parent)
	if err != nil {
		return nil, err
	}

	for i := range departments {
		departments[i].Departments, err = departmentAPI.getTree(&departments[i])
		if err != nil {
			return nil, err
		}
	}

	return departments, nil
}

func (departmentAPI *DepartmentAPI) getList(parent *Department) ([]Department, error) {
	departments, err := departmentAPI.db.GetDepartments(parent)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(departments); i++ {
		children, err := departmentAPI.getList(&departments[i])
		if err != nil {
			return nil, err
		}

		// inserting children into slice after parent
		slicingIndex := i + 1
		departments = append(departments[:slicingIndex], append(children, departments[slicingIndex:]...)...)
		i += len(children)
	}

	return departments, nil
}

func (departmentAPI *DepartmentAPI) save(department *Department) error {
	return departmentAPI.db.SaveDepartment(department)
}

func (departmentAPI *DepartmentAPI) delete(id int) error {
	department, err := departmentAPI.db.GetDepartment(id)
	if err != nil {
		return err
	}

	isReferenced, err := departmentAPI.db.IsDepartmentReferenced(department)
	if err != nil {
		return err
	} else if isReferenced {
		return errItemInUse
	}

	err = departmentAPI.db.DeleteDepartment(department)
	if err != nil {
		return err
	}

	return nil
}
