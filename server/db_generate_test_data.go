package server

import (
	"math/rand"
	"time"

	"gopkg.in/gorp.v1"
)

func (db *DB) GenerateTestData(testPwdHash []byte) error {
	trans, err := db.dbMap.Begin()
	if err != nil {
		return err
	}

	// warning:
	// specified id in objects to insert will be overridden by auto-inc values in db.
	// here, they reflect the indexes of the slices, to referenced them more easy in further inserts/updates.

	departments := []interface{}{
		&Department{0, nil, "Dept. 1", 0, "", nil},
		&Department{1, nil, "Dept. 1a", 0, "", nil},
		&Department{2, nil, "Dept. 1a1", 0, "", nil},
		&Department{3, nil, "Dept. 1b", 0, "", nil}}
	if err = db.insertTestData(trans, &departments); err != nil {
		return err
	}

	users := []interface{}{
		&User{0, intPtr(RoleAdmin), &departments[0].(*Department).ID, "admin", testPwdHash, "Nny", "C.", "admin@example.com", 0, "", ""},
		&User{1, intPtr(RoleManager), &departments[1].(*Department).ID, "manager", testPwdHash, "Zim", "Irken", "manager@example.com", 0, "", ""},
		&User{2, intPtr(RoleDeptManager), &departments[1].(*Department).ID, "deptman", testPwdHash, "Tak", "Irken", "deptman@example.com", 0, "", ""},
		&User{3, intPtr(RoleUser), &departments[3].(*Department).ID, "user", testPwdHash, "GIR", "Dimmwitted", "user@example.com", 0, "", ""}}
	if err = db.insertTestData(trans, &users); err != nil {
		return err
	}

	projectCategories := []interface{}{
		&ProjectCategory{0, nil, "AA", "AA", "Category 1", 0, "", nil},
		&ProjectCategory{1, nil, "01", "AA01", "Category 1>1", 0, "", nil},
		&ProjectCategory{2, nil, "02", "AA02", "Category 1>2", 0, "", nil},
		&ProjectCategory{3, nil, "001", "AA02001", "Category 1>2>1", 0, "", nil},
		&ProjectCategory{4, nil, "03", "AA03", "Category 1>3", 0, "", nil},
		&ProjectCategory{5, nil, "BB", "BB", "Category 2", 0, "", nil},
		&ProjectCategory{6, nil, "01", "BB01", "Category 2>1", 0, "", nil},
		&ProjectCategory{7, nil, "02", "BB02", "Category 2>2", 0, "", nil}}
	if err = db.insertTestData(trans, &projectCategories); err != nil {
		return err
	}

	projects := []interface{}{
		&Project{0, projectCategories[1].(*ProjectCategory).ID, "01", "AA/01", &users[0].(*User).ID, &users[1].(*User).ID, "Project 1", "", 0, nil, nil, nil},
		&Project{1, projectCategories[3].(*ProjectCategory).ID, "01", "AA02001/01", &users[0].(*User).ID, &users[2].(*User).ID, "Project 2", "", 0, nil, nil, nil},
		&Project{2, projectCategories[6].(*ProjectCategory).ID, "01", "BB01/01", &users[1].(*User).ID, &users[0].(*User).ID, "Project 3", "", 0, nil, nil, nil},
		&Project{3, projectCategories[6].(*ProjectCategory).ID, "02", "BB01/02", &users[1].(*User).ID, &users[3].(*User).ID, "Project 4", "", 0, nil, nil, nil}}
	if err = db.insertTestData(trans, &projects); err != nil {
		return err
	}

	activityTypes := []interface{}{
		&ActivityType{0, "Activity Type 1", 0},
		&ActivityType{1, "Activity Type 2", 0},
		&ActivityType{2, "Activity Type 3", 0}}
	if err = db.insertTestData(trans, &activityTypes); err != nil {
		return err
	}

	projectActivityTypes := []interface{}{
		&ProjectActivityType{projects[0].(*Project).ID, activityTypes[0].(*ActivityType).ID},
		&ProjectActivityType{projects[0].(*Project).ID, activityTypes[1].(*ActivityType).ID},
		&ProjectActivityType{projects[1].(*Project).ID, activityTypes[0].(*ActivityType).ID},
		&ProjectActivityType{projects[1].(*Project).ID, activityTypes[1].(*ActivityType).ID},
		&ProjectActivityType{projects[1].(*Project).ID, activityTypes[2].(*ActivityType).ID},
		&ProjectActivityType{projects[2].(*Project).ID, activityTypes[1].(*ActivityType).ID},
		&ProjectActivityType{projects[2].(*Project).ID, activityTypes[2].(*ActivityType).ID},
		&ProjectActivityType{projects[3].(*Project).ID, activityTypes[0].(*ActivityType).ID},
		&ProjectActivityType{projects[3].(*Project).ID, activityTypes[2].(*ActivityType).ID}}
	if err = db.insertTestData(trans, &projectActivityTypes); err != nil {
		return err
	}

	projectDepartments := []interface{}{
		&ProjectDepartment{projects[0].(*Project).ID, departments[0].(*Department).ID},
		&ProjectDepartment{projects[1].(*Project).ID, departments[1].(*Department).ID},
		&ProjectDepartment{projects[2].(*Project).ID, departments[2].(*Department).ID},
		&ProjectDepartment{projects[3].(*Project).ID, departments[1].(*Department).ID},
		&ProjectDepartment{projects[3].(*Project).ID, departments[2].(*Department).ID}}
	if err = db.insertTestData(trans, &projectDepartments); err != nil {
		return err
	}

	projectUsers := []interface{}{
		&ProjectUser{projects[0].(*Project).ID, users[0].(*User).ID},
		&ProjectUser{projects[1].(*Project).ID, users[1].(*User).ID},
		&ProjectUser{projects[2].(*Project).ID, users[2].(*User).ID},
		&ProjectUser{projects[3].(*Project).ID, users[0].(*User).ID},
		&ProjectUser{projects[3].(*Project).ID, users[1].(*User).ID},
		&ProjectUser{projects[3].(*Project).ID, users[2].(*User).ID}}
	if err = db.insertTestData(trans, &projectUsers); err != nil {
		return err
	}

	departments[1].(*Department).ParentID = &departments[0].(*Department).ID
	departments[2].(*Department).ParentID = &departments[1].(*Department).ID
	departments[3].(*Department).ParentID = &departments[0].(*Department).ID
	if err = db.updateTestData(trans, &departments); err != nil {
		return err
	}

	projectCategories[1].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[2].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[3].(*ProjectCategory).ParentID = &projectCategories[1].(*ProjectCategory).ID
	projectCategories[4].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[6].(*ProjectCategory).ParentID = &projectCategories[5].(*ProjectCategory).ID
	projectCategories[7].(*ProjectCategory).ParentID = &projectCategories[5].(*ProjectCategory).ID
	if err = db.updateTestData(trans, &projectCategories); err != nil {
		return err
	}

	activities := db.generateActivityData(users, projects, departments, projectUsers, projectDepartments, projectActivityTypes)
	if err = db.insertTestData(trans, &activities); err != nil {
		return err
	}

	return trans.Commit()
}

func (db *DB) generateActivityData(users []interface{}, projects []interface{}, departments []interface{}, projectUsers []interface{}, projectDepartments []interface{}, projectActivityTypes []interface{}) []interface{} {
	currentYear := time.Now().Year()
	date := time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.UTC)
	var activities []interface{}

	mapUserProjectIDs := db.generateUserProjectMap(users, projects, departments, projectUsers, projectDepartments)
	mapProjectActivityTypeIDs := db.generateProjectActivityTypeMap(projectActivityTypes)

	for date.Year() == currentYear {
		for userID, userProjectIDs := range mapUserProjectIDs {
			activityCount := db.determineActivityCountPerDay(mapProjectActivityTypeIDs, userProjectIDs)
			var mapUsedProjectActivityTypes = make(map[int][]int)

			for i := 0; i < activityCount; i++ {
				duration := rand.Intn(1440 / activityCount)

				projectID, activityTypeID := db.generateIDsProjectActivityType(userProjectIDs, mapProjectActivityTypeIDs, mapUsedProjectActivityTypes)
				mapUsedProjectActivityTypes[projectID] = append(mapUsedProjectActivityTypes[projectID], activityTypeID)

				activity := &Activity{-1, date, userID, projectID, activityTypeID, duration, 1}
				activities = append(activities, activity)
			}

		}
		date = date.AddDate(0, 0, 1)
	}

	return activities
}

func (db *DB) generateUserProjectMap(users []interface{}, projects []interface{}, departments []interface{}, projectUsers []interface{}, projectDepartments []interface{}) map[int][]int {
	var mapUserProjectIDs = make(map[int][]int)

	for i := range projectUsers {
		projectUser := projectUsers[i].(*ProjectUser)
		mapUserProjectIDs[projectUser.UserID] = append(mapUserProjectIDs[projectUser.UserID], projectUser.ProjectID)
	}

	db.addManagerRespToUserProjectMap(projects, mapUserProjectIDs)
	db.addProjectIDsOfDepartment(users, departments, projectDepartments, mapUserProjectIDs)

	return mapUserProjectIDs
}

func (db *DB) addManagerRespToUserProjectMap(projects []interface{}, mapUserProjectIDs map[int][]int) {

	for _, value := range projects {
		project := value.(*Project)

		addResponsible := true
		addManager := true

		for _, projectID := range mapUserProjectIDs[*project.ResponsibleUserID] {
			if projectID == project.ID {
				addResponsible = false
				break
			}
		}

		for _, projectID := range mapUserProjectIDs[*project.ManagerUserID] {
			if projectID == project.ID {
				addManager = false
				break
			}
		}

		if addResponsible {
			mapUserProjectIDs[*project.ResponsibleUserID] = append(mapUserProjectIDs[*project.ResponsibleUserID], project.ID)
		}
		if addManager {
			mapUserProjectIDs[*project.ManagerUserID] = append(mapUserProjectIDs[*project.ManagerUserID], project.ID)
		}
	}
}

func (db *DB) addProjectIDsOfDepartment(users []interface{}, depts []interface{}, projectDepts []interface{}, mapUserProjectIDs map[int][]int) {
	for _, u := range users {
		user := u.(*User)
		projectIDs := db.getAllProjectIDsOfDepartment(*user.DepartmentID, depts, projectDepts)

		for _, projectID := range projectIDs {
			add := true

			for _, userProjectID := range mapUserProjectIDs[user.ID] {
				if projectID == userProjectID {
					add = false
					break
				}
			}

			if add {
				mapUserProjectIDs[user.ID] = append(mapUserProjectIDs[user.ID], projectID)
			}
		}
	}
}

func (db *DB) getAllProjectIDsOfDepartment(deptID int, depts []interface{}, projectDepts []interface{}) []int {
	projectIDs := db.getProjectIDsOfDepartment(deptID, projectDepts)
	dept := db.getDepartment(deptID, depts)

	for dept.ParentID != nil {
		dept = db.getDepartment(*dept.ParentID, depts)
		projectIDs = append(projectIDs, db.getProjectIDsOfDepartment(dept.ID, projectDepts)...)
	}

	return projectIDs
}

func (db *DB) getProjectIDsOfDepartment(deptID int, projectDepts []interface{}) []int {
	var projectIDs []int

	for _, pd := range projectDepts {
		projectDept := pd.(*ProjectDepartment)
		if projectDept.DepartmentID == deptID {
			projectIDs = append(projectIDs, projectDept.ProjectID)
		}
	}

	return projectIDs
}

func (db *DB) getDepartment(deptID int, depts []interface{}) Department {
	var dept Department

	for _, d := range depts {
		dept = *d.(*Department)
		if dept.ID == deptID {
			break
		}
	}

	return dept
}

func (db *DB) generateProjectActivityTypeMap(projectActivityTypes []interface{}) map[int][]int {
	var mapProjectActivityTypeIDs = make(map[int][]int)

	for i := range projectActivityTypes {
		projectActivityType := projectActivityTypes[i].(*ProjectActivityType)
		mapProjectActivityTypeIDs[projectActivityType.ProjectID] = append(mapProjectActivityTypeIDs[projectActivityType.ProjectID], projectActivityType.ActivityTypeID)
	}

	return mapProjectActivityTypeIDs
}

func (db *DB) determineActivityCountPerDay(mapProjectActivityTypeIDs map[int][]int, userProjectIDs []int) int {
	activityCount := 0

	for k, v := range mapProjectActivityTypeIDs {
		add := false
		for _, v := range userProjectIDs {
			if k == v {
				add = true
				break
			}
		}
		if !add {
			continue
		}

		// don't add more than 4 activities per day
		activityCount += len(v)
		if activityCount > 4 {
			activityCount = 4
			break
		}
	}

	return activityCount
}

func (db *DB) generateIDsProjectActivityType(userProjectIDs []int, mapProjectActivityTypeIDs map[int][]int, mapUsedProjectActivityTypes map[int][]int) (int, int) {
	var projectID, activityTypeID int
	generateIDs := true

	for generateIDs {
		projectID = userProjectIDs[rand.Intn(len(userProjectIDs))]
		activityTypeID = mapProjectActivityTypeIDs[projectID][rand.Intn(len(mapProjectActivityTypeIDs[projectID]))]

		usedActivtyTypes, ok := mapUsedProjectActivityTypes[projectID]
		if !ok {
			break
		}

		for _, v := range usedActivtyTypes {
			if v == activityTypeID {
				generateIDs = true
				break
			} else {
				generateIDs = false
			}
		}
	}

	return projectID, activityTypeID
}

func (db *DB) insertTestData(trans *gorp.Transaction, data *[]interface{}) error {
	err := trans.Insert(*data...)
	if err != nil {
		trans.Rollback()
	}
	return err
}

func (db *DB) updateTestData(trans *gorp.Transaction, data *[]interface{}) error {
	_, err := trans.Update(*data...)
	if err != nil {
		trans.Rollback()
	}
	return err
}
