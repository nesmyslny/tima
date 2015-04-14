package server

import "gopkg.in/gorp.v1"

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
		&Project{0, projectCategories[1].(*ProjectCategory).ID, "01", "AA/01", &users[0].(*User).ID, &users[1].(*User).ID, "Project 1", 0, nil, nil, nil},
		&Project{1, projectCategories[3].(*ProjectCategory).ID, "01", "AA02001/01", &users[0].(*User).ID, &users[2].(*User).ID, "Project 2", 0, nil, nil, nil},
		&Project{2, projectCategories[6].(*ProjectCategory).ID, "01", "BB01/01", &users[1].(*User).ID, &users[0].(*User).ID, "Project 3", 0, nil, nil, nil},
		&Project{3, projectCategories[6].(*ProjectCategory).ID, "02", "BB01/02", &users[1].(*User).ID, &users[3].(*User).ID, "Project 4", 0, nil, nil, nil}}
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

	// todo: generate activity data (maybe random data for the current month?)

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

	return trans.Commit()
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
