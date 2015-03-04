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

	users := []interface{}{
		&User{0, intPtr(RoleAdmin), "admin", testPwdHash, "Nny", "C.", "admin@example.com", "", ""},
		&User{1, intPtr(RoleManager), "manager", testPwdHash, "Zim", "Irken", "manager@example.com", "", ""},
		&User{2, intPtr(RoleUser), "user", testPwdHash, "GIR", "Dimmwitted", "user@example.com", "", ""}}
	if err = db.insertTestData(trans, users); err != nil {
		return err
	}

	projectCategories := []interface{}{
		&ProjectCategory{0, nil, "AA", "AA", "Category 1", "", nil},
		&ProjectCategory{1, nil, "01", "AA01", "Category 1>1", "", nil},
		&ProjectCategory{2, nil, "02", "AA02", "Category 1>2", "", nil},
		&ProjectCategory{3, nil, "001", "AA02001", "Category 1>2>1", "", nil},
		&ProjectCategory{4, nil, "03", "AA03", "Category 1>3", "", nil},
		&ProjectCategory{5, nil, "BB", "BB", "Category 2", "", nil},
		&ProjectCategory{6, nil, "01", "BB01", "Category 2>1", "", nil},
		&ProjectCategory{7, nil, "02", "BB02", "Category 2>2", "", nil}}
	if err = db.insertTestData(trans, projectCategories); err != nil {
		return err
	}

	projects := []interface{}{
		&Project{0, projectCategories[1].(*ProjectCategory).ID, "01", "AA/01", &users[0].(*User).ID, &users[1].(*User).ID, "Project 1", nil},
		&Project{1, projectCategories[3].(*ProjectCategory).ID, "01", "AA02001/01", &users[0].(*User).ID, &users[2].(*User).ID, "Project 2", nil},
		&Project{2, projectCategories[6].(*ProjectCategory).ID, "01", "BB01/01", &users[1].(*User).ID, &users[0].(*User).ID, "Project 3", nil},
		&Project{3, projectCategories[7].(*ProjectCategory).ID, "02", "BB02/02", &users[1].(*User).ID, &users[2].(*User).ID, "Project 4", nil}}
	if err = db.insertTestData(trans, projects); err != nil {
		return err
	}

	activityTypes := []interface{}{
		&ActivityType{0, "Activity Type 1"},
		&ActivityType{1, "Activity Type 2"},
		&ActivityType{2, "Activity Type 3"}}
	if err = db.insertTestData(trans, activityTypes); err != nil {
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
	if err = db.insertTestData(trans, projectActivityTypes); err != nil {
		return err
	}

	// todo: generate activity data (maybe random data for the current month?)

	projectCategories[1].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[2].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[3].(*ProjectCategory).ParentID = &projectCategories[1].(*ProjectCategory).ID
	projectCategories[4].(*ProjectCategory).ParentID = &projectCategories[0].(*ProjectCategory).ID
	projectCategories[6].(*ProjectCategory).ParentID = &projectCategories[5].(*ProjectCategory).ID
	projectCategories[7].(*ProjectCategory).ParentID = &projectCategories[5].(*ProjectCategory).ID
	if err = db.updateTestData(trans, projectCategories); err != nil {
		return err
	}

	return trans.Commit()
}

func (db *DB) insertTestData(trans *gorp.Transaction, data []interface{}) error {
	err := trans.Insert(data...)
	if err != nil {
		trans.Rollback()
	}
	return err
}

func (db *DB) updateTestData(trans *gorp.Transaction, data []interface{}) error {
	_, err := trans.Update(data...)
	if err != nil {
		trans.Rollback()
	}
	return err
}
