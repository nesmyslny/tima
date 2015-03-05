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

func jsonResult(boolResult bool, stringResult string, intResult int) (interface{}, *HandlerError) {
	return JsonResult{boolResult, stringResult, intResult}, nil
}

func jsonResultBool(boolResult bool) (interface{}, *HandlerError) {
	return JsonResult{BoolResult: boolResult}, nil
}

func jsonResultString(stringResult string) (interface{}, *HandlerError) {
	return JsonResult{StringResult: stringResult}, nil
}

func jsonResultInt(intResult int) (interface{}, *HandlerError) {
	return JsonResult{IntResult: intResult}, nil
}

func intPtr(i int) *int {
	return &i
}
