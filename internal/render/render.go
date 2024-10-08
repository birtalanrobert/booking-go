package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/birtalanrobert/bookings-go/internal/config"
	"github.com/birtalanrobert/bookings-go/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{
	"humanDate": HumanDate,
	"formatDate": FormatDate,
	"iterate": Iterate,
	"add": Add,
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in yyyy-mm-dd
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate returns a slice of items, starting at 1, going to count
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}

	return items
}

func Add(a int, b int) int {
	return a + b
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateDate, r *http.Request) *models.TemplateDate {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateDate) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Println("Could not get template from template cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.go.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.go.tmpl", pathToTemplates))
		
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.go.tmpl", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}