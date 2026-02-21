package middlewares

func GetMiddleware() *Middleware {
	return globalMiddleware
}
func GetAnyRoleMiddleware() *AnyRoleMiddleware {
	return globalMiddleware
}
