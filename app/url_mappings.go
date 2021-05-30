package app

import (
	"github.com/harshasavanth/users-api/controllers/ping"
	"github.com/harshasavanth/users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/users", users.UsersController.Create)
	router.GET("/users/:user_id", users.UsersController.IsAuthorized(users.UsersController.Get))
	router.GET("/users/byemail/:email", users.UsersController.GetByEmail)
	router.PUT("users/:user_id", users.UsersController.IsAuthorized(users.UsersController.Update))
	router.GET("users/verifyemail/:id", users.UsersController.VerifyEmail)
	//router.PATCH("users/:user_id", users.Update)
	router.DELETE("users/:user_id", users.UsersController.IsAuthorized(users.UsersController.Delete))
	//router.GET("/internal/users/search", users.Search)
	//router.POST("/users/login", users.Login)
}
