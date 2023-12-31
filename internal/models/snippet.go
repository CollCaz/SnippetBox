package models

import (
	"database/sql"
	"errors"
	"time"
)

// A Struct for the data in each snippet in our database,
// The fields of the struct corrosponde to the fields in our MySQL snippets
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// Struct wrapping sql.DB method
type SnippetModel struct {
	DB *sql.DB
}

// snippetModel method to insert a new snippet to the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

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

// snippetModel method to get a snippet base on ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// snippetModel method to get the latest 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
