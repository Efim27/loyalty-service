package handlers

import "loyalty-service/internal/middleware"

func SetupRoutes(server Server) {
	middlewareLoginRequired := middleware.NewLoginRequired(server.Config.Secret)
	api := server.App.Group("/api")

	apiUser := api.Group("/user")
	apiUser.Post("/register", server.userRegister)
	apiUser.Post("/login", server.userLogin)
	apiUser.Get("/withdrawals", middlewareLoginRequired, server.withdrawalList)

	apiOrder := apiUser.Group("/orders")
	apiOrder.Post("/", middlewareLoginRequired, server.orderNew)
	apiOrder.Get("/", middlewareLoginRequired, server.orderList)

	apiBalance := apiUser.Group("/balance")
	apiBalance.Get("/", middlewareLoginRequired, server.balanceGet)
	apiBalance.Post("/withdraw", middlewareLoginRequired, server.withdrawalNew)
}
