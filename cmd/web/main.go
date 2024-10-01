package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/ericmccorkle/lets-go/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to hold application-wide dependencies for
// the app.
type application struct {
	logger *slog.Logger
	// This allows us to make the SnippetModel object available to handlers
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Define a new command-line flag with name 'addr', a default value
	// of ":4000" and some short help text explaining what the flag controls.
	// The value of the flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// A command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Parse the command line flag. It reads the value and assigns it
	// to the addr variable. If any erros are encountered, the application
	// will be terminated.
	flag.Parse()

	// Initialize a new structured logger with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// defer so that connection pool is closed before main() exits.
	defer db.Close()

	// Initialize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// New instance of the application struct, containing dependencies
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	// Log any error message returned by ListenAndServe() at Error severity.
	// Then call os.Exit(1) to terminate the app with exit code 1.
	// Slog has no equivalent of log.Fatal(). This is the closest we can get.
	logger.Error(err.Error())
	os.Exit(1)
}

// openDB() wraps slq.Open() and returns a sql.DB connection pool
// for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// sql.Open() doesn't actually create any connections, but rather it
	// initializes the pool for future use. Connections are established
	// lazily when needed for the first time.
	// On set up, we need to use db.Ping() to create a connection and
	// check for any errors.
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
