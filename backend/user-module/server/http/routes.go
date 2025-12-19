package http

// RegisterRoutes registers all HTTP routes for the User module
func RegisterRoutes(he *UserHttpExtension) {
	he.Init()
}
