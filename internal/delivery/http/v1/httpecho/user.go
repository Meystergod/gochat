package httpecho

import (
	"github.com/Meystergod/gochat/internal/controller"

	"github.com/labstack/echo/v4"
)

func SetUserApiRoutes(e *echo.Echo, userController *controller.UserController) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/user", userController.CreateUser)
		v1.GET("/user/:id", userController.GetUserInfo)
		v1.GET("/users", userController.GetAllUsersInfo)
		v1.PUT("/user/:id", userController.UpdateUserInfo)
		v1.DELETE("/user/:id", userController.DeleteUserAccount)
	}
}
