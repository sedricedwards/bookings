package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sedricedwards/bookings/pkg/config"
	"github.com/sedricedwards/bookings/pkg/handlers"
	"github.com/sedricedwards/bookings/pkg/render"
)

const portNumber = 8080

var app config.AppConfig
var session *scs.SessionManager

// main is the entry point for the application
func main() {

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", portNumber),
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
