package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/birtalanrobert/bookings-go/internal/config"
	"github.com/birtalanrobert/bookings-go/internal/driver"
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
	db, err := run();	

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()
	defer close(app.MailChan)
	listenForMail()

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

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=robertbirtalan password=")

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to database successfully!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelper(&app)

	return db, nil
}