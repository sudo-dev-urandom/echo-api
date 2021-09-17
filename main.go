package main

import (
	"goecho/controllers"
	"goecho/core"
	_ "goecho/docs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	defer core.App.Close()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Pre(middleware.Rewrite(map[string]string{
		"/api/*": "/$1",
	}))

	e.GET("/docs/*", echoSwagger.WrapHandler)

	api := e.Group("/v1")
	{

		users := api.Group("/users")
		users.GET("", controllers.UserList)
		users.POST("/create", controllers.UserCreate)
		users.PATCH("/update/:id", controllers.UserUpdate)
		users.DELETE("/delete/:id", controllers.UserDelete)

	}

	e.Logger.Fatal(e.Start(":" + core.App.Port))
	os.Exit(0)

}
