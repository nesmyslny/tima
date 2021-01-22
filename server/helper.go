package server

func AuthorizeUser(context *HandlerContext) (bool, error) {
	return *context.User.Role >= RoleUser, nil
}

func AuthorizeDeptManager(context *HandlerContext) (bool, error) {
	return *context.User.Role >= RoleDeptManager, nil
}

func AuthorizeManager(context *HandlerContext) (bool, error) {
	return *context.User.Role >= RoleManager, nil
}

func AuthorizeAdmin(context *HandlerContext) (bool, error) {
	return *context.User.Role >= RoleAdmin, nil
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func diffInt(x, y []int) []int {
	var diff []int

	for _, i := range x {
		add := true
		for _, j := range y {
			if i == j {
				add = false
				break
			}
		}
		if add {
			diff = append(diff, i)
		}
	}
	return diff
}

func createSqlArgString(count int) string {
	var str string
	for i := 0; i < count; i++ {
		if i > 0 {
			str += ", "
		}
		str += "?"
	}
	return str
}

func sliceIntToInterface(in []int) []interface{} {
	out := make([]interface{}, len(in))
	for idx, val := range in {
		out[idx] = val
	}
	return out
}

// func getProjectIDsAsInterfaces(projects []Project) []interface{} {
// 	var IDs []interface{}
// 	for _, p := range projects {
// 		IDs = append(IDs, p.ID)
// 	}
//
// 	return IDs
// }

func getProjectIDs(projects []Project) []int {
	var IDs []int
	for _, p := range projects {
		IDs = append(IDs, p.ID)
	}

	return IDs
}
