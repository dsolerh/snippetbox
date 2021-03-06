package mysql

import (
	"database/sql"
	"io/ioutil"
	"strings"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	// Establish a sql.DB connection pool for our test database. Because our
	// setup and teardown scripts contains multiple SQL statements, we need
	// to use the `multiStatements=true` parameter in our DSN. This instructs
	// our MySQL database driver to support executing multiple SQL statements
	// in one db.Exec()` call.
	db, err := sql.Open("mysql", "web:pass@tcp(localhost:3306)/snippetbox_test?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from file and execute the statements.
	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	stmts := strings.SplitAfter(string(script), ";")
	for _, stmt := range stmts {
		if stmt == "" {
			continue
		}
		_, err = db.Exec(stmt)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Return the connection pool and an anonymous function which reads and
	// executes the teardown script, and closes the connection pool. We can
	// assign this anonymous function and call it later once our test has
	// completed.
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		stmts := strings.SplitAfter(string(script), ";")
		for _, stmt := range stmts {
			if stmt == "" {
				continue
			}
			_, err = db.Exec(stmt)
			if err != nil {
				t.Fatal(err)
			}
		}
		db.Close()
	}
}
