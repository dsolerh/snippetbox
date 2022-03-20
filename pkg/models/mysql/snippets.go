package mysql

import (
	"database/sql"

	"dsolerh.projects/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB

	// prepared statements
	InsertStmt *sql.Stmt
	GetStmt    *sql.Stmt
	LatestStmt *sql.Stmt
}

func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {
	insertStmt, err := db.Prepare(`INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`)
	if err != nil {
		return nil, err
	}
	getStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`)
	if err != nil {
		return nil, err
	}
	latestStmt, err := db.Prepare(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}

	return &SnippetModel{db, insertStmt, getStmt, latestStmt}, nil
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	result, err := m.InsertStmt.Exec(title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	row := m.GetStmt.QueryRow(id)

	// create an empty Snippet
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// This will return the latest Snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := m.LatestStmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	// iterate over the rows
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append to the slice of snippets.
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

func (m *SnippetModel) ClosePrepared() {
	m.InsertStmt.Close()
	m.GetStmt.Close()
	m.LatestStmt.Close()
}
