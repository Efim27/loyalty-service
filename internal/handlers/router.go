package handlers

import "loyalty-service/internal/middleware"

func SetupRoutes(server Server) {
	middlewareLoginRequired := middleware.NewLoginRequired(server.Config.Secret)
	api := server.App.Group("/api")

	apiUser := api.Group("/user")
	apiUser.Post("/register", server.userRegister)
	apiUser.Post("/login", server.userLogin)

	apiOrder := apiUser.Group("/orders")
	apiOrder.Post("/", middlewareLoginRequired, server.orderNew)

}
