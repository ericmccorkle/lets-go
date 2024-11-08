package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/ericmccorkle/lets-go/snippetbox/internal/models"
)

// templateData type acts as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

// humanDate returns a human formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	return t.Format("01 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This is a lookup between the names of our custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// A new map to act as the cache.
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths that match the given pattern.
	// This gives us a slice of all filepaths for our application 'page' templates
	// Ex: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the file name (like 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you call the ParseFiles()
		// method. This means we have to use template.New() to create an empty template set,
		// use the Funcs() method to register the template.FunctMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on this template set to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add template set to the map, using the name of the page as the key.
		cache[name] = ts
	}
	return cache, nil
}
