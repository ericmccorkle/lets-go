package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	// Initialize a slice to contain the paths to the files.
	// The file containing the base template must be the first file
	// in the slice.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Use template.ParseFiles() func to read the template file into a
	// template set. If an error, log the detailed error message, use the http.Error()
	// func to send an Internal Server Error response to the user,
	// and then return from the handler so no subsequent code is executed.
	// the ... after files passes the contents of the slice as variadic arguments.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Then use the Execute() method on the template set to write the template
	// content as response body. Last parameter to Execute() represents any
	// dynamic data that we want to pass in.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	testId, err := strconv.Atoi(r.PathValue("testId"))
	if err != nil {
		fmt.Println("Error with testId")
	}
	runId, err := strconv.Atoi(r.PathValue("runId"))
	if err != nil {
		fmt.Println("Error with runId")
	}
	fmt.Println("testId", testId)
	fmt.Println("runId", runId)
}
