package server

func AuthorizeUser(context *HandlerContext) (bool, error) {
	return *context.User.Role >= RoleUser, nil
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
