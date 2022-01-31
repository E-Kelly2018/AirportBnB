package main

import (
	"encoding/gob"
	"fmt"
	"github.com/E-Kelly2018/AirportBnB/internal/config"
	"github.com/E-Kelly2018/AirportBnB/internal/handlers"
	"github.com/E-Kelly2018/AirportBnB/internal/models"
	"github.com/E-Kelly2018/AirportBnB/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

const port = ":8080"

func main() {
	gob.Register(models.Reservation{})
	//Chnage to true when in production
	app.InProduction = false

	//Set up session timeout
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on %s", port))
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
