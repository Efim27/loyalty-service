package handlers

func SetupRoutes(server Server) {
	//api := app.Group("/api", middleware.AuthReq())
	api := server.App.Group("/api")

	apiUser := api.Group("/user")
	apiUser.Post("/register", server.userRegister)
	apiUser.Post("/login", server.userLogin)
}
