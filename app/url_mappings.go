package app

import (
	"github.com/harshasavanth/users-api/controllers/ping"
	"github.com/harshasavanth/users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.POST("/register", users.UsersController.Create)
	router.GET("/users/:user_id", users.UsersController.IsAuthorized(users.UsersController.Get))
	router.PUT("/update/:user_id", users.UsersController.IsAuthorized(users.UsersController.Update))
	router.DELETE("/delete/:user_id", users.UsersController.IsAuthorized(users.UsersController.Delete))
	router.GET("/users/getbyemail/:email", users.UsersController.GetByEmail)
	router.GET("/users/verifyemail/:id", users.UsersController.VerifyEmail)

}
