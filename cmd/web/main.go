package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Define a new command-line flag with name 'addr', a default value
	// of ":4000" and some short help text explaining what the flag controls.
	// The value of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command line flag. It reads the value and assigns it
	// to the addr variable. If any erros are encountered, the application
	// will be terminated.
	flag.Parse()

	// Initialize a new structured logger with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	// Create a file server that serves files out of the "./ui/static" dir.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// For matching paths, we strip the "/static" prefix before the request
	// reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Testing a path value in middle of path
	mux.HandleFunc("GET /tests/{testId}/runs/{runId}", testFunc)

	logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)

	// Log any error message returned by ListenAndServe() at Error severity.
	// Then call os.Exit(1) to terminate the app with exit code 1.
	// Slog has no equivalent of log.Fatal(). This is the closest we can get.
	logger.Error(err.Error())
	os.Exit(1)
}
