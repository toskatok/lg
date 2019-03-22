package actions

import (
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App creates new instance of Echo and configures it
func App() *echo.Echo {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	// Validator
	app.Validator = &DefaultValidator{validator.New()}

	// Routes
	app.GET("/about", AboutHandler)
	api := app.Group("/api")
	{
		ih := NewInstancesHandler()
		api.POST("/instances", ih.Create)
		api.GET("/instances", ih.List)
		api.DELETE("/instances/:instance_id", ih.Destroy)
		api.GET("/instances/:instance_id", ih.Show)
	}

	return app
}
