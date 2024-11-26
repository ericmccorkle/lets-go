package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server
// Error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page
	// name. If no entry exists in the cache with the provided name, then error.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exists", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write template to buffer, instead of straight to the http.ResponseWriter.
	// If an error, call serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If template is written tot he buffer without any errors, we are safe to go
	// ahead and write the HTTP status code to the http.ResponseWriter
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	buf.WriteTo(w)
}

// newTemplateData returns a pointer to a templateData struct initialized with the current year
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

// decodePostForm takes a pointer to a http.Request and a dst (a target "destination")
// that we want to decode form data into
// We use this function as a way to avoid a nil-pointer error if target destination
// is invalid
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance. Pass in target destination as first parameter
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// if we use an invalid target destination, Decode() method will return an error
		// with the type *form.InvalidDecoderError.
		// We use errors.As() to check for this and raise a panic instead of return the error
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
