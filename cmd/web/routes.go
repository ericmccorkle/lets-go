package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes() method returns a handler containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Create a standard middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	// Pass the mux ServeMux instance as the "next" parameter to the commonHeaders
	// middleware. Because commonHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
