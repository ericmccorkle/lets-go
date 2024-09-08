package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// The SnippetModel wraps a sql.DB connection pool
// If helpful, think of "model" as a service layer or data access layer
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Returns a sql.Result type, containing info on what happened at execution
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// QueryRow() returns a pointer to sql.Row object that holds the result from the db
	// Note: errors from DB.QueryRow() are deferred until Scan() is called, so it is
	// possibel to instead write something like:
	// err := m.DB.QueryRow("SELECT...", id).Scan(&s.ID, &s.Title,...)
	row := m.DB.QueryRow(stmt, id)

	var s Snippet
	// row.Scan() copies the values from each field in sql.Row to their corresponding fields
	// in the Snippet struct. The arguments are pointers to the place you want to copy
	// the data into.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If no rows returned, row.Scan() returns a sql.ErrNoRows error.
		// We use the errors.Is() functio nto check for that specific error.
		// Then return our own ErrNoRecord instead.
		// We create our own ErrNoRecord error so that we encapsulate the model
		// completely. Our handlers don't need to be concerned with the underlying
		// datastore or reliant on datastore-specific error (like sql.ErrNoRows) for
		// its behavior.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Ensure the sql.Rows resultset is properly closed before the method returns.
	// Should come after the error check. Otherwise, if Query() returns an error,
	// you'll get a panic trying to lcose a nil resultset.
	defer rows.Close()

	var snippets []Snippet
	// Use rows.Next to itereate through the rows in the resultset.
	// THis prepares the first row ot be acted on by the rows.Sna() method.
	// If iteration over all rows completes then the resultset automatically
	// closes itself and frees up the underlying database connection.
	for rows.Next() {
		// Craete a pointer to a new zeroed Snippet struct
		var s Snippet
		// Use rows. Scan to copy the values from each field in the row to the new
		// Snippet object that we created.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// After the loop has completed, call rows.Err() to retrieve any
	// error encountered during the iteration. Don't assume a successful
	// iteration was completed over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
