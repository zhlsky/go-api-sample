package main

import (
	. "echo-sample/handlers"
	"github.com/labstack/echo"
)

func (app *app) initRouter() {
	v1 := app.Group("/api/v1")
	{
		v1.GET("/employee", echo.HandlerFunc(GetEmployees))
		v1.GET("/employee/:id", echo.HandlerFunc(GetEmployee))
		v1.POST("/employee", echo.HandlerFunc(CreateEmployee))
		v1.PATCH("/employee/:id", echo.HandlerFunc(UpdateEmployee))
		v1.DELETE("/employee/:id", echo.HandlerFunc(DeleteEmployee))
	}
}
