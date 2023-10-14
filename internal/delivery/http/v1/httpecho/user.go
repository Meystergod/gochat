package httpecho

import (
	"github.com/labstack/echo/v4"
)

func SetUserApiRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/user", nil)
		v1.GET("/user/:id", nil)
		v1.GET("/users", nil)
		v1.PUT("/user/:id", nil)
		v1.DELETE("/user/:id", nil)
	}
}
