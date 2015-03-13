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
