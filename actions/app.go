package actions

import (
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/gobuffalo/x/sessions"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var validate *validator.Validate
var socket *socketio.Server

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_lg_session",
		})

		// If no content type is sent by the client
		// the application/json will be set, otherwise the client's
		// content type will be used.
		app.Use(contenttype.Add("application/json"))

		// validator
		validate = validator.New()

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// socket.io initiation
		sio, err := socketio.NewServer(nil)
		if err != nil {
			log.Fatal(err)
		}
		socket = sio // moves it into global scope

		// Routes
		// swagger ui
		app.ServeFiles("/swagger", http.Dir("swagger"))
		app.GET("/about", AboutHandler)
		app.Mount("/socket.io/", sio) // handles the socket io
		api := app.Group("/api")
		{
			api.Resource("/instances", InstancesResource{})
		}
	}

	return app
}
