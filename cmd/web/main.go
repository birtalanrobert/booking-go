package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/birtalanrobert/bookings-go/internal/config"
	"github.com/birtalanrobert/bookings-go/internal/handlers"
	"github.com/birtalanrobert/bookings-go/internal/helpers"
	"github.com/birtalanrobert/bookings-go/internal/models"
	"github.com/birtalanrobert/bookings-go/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run();	

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelper(&app)

	return nil
}