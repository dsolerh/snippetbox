package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	// my package for snippet related functionalities
	"dsolerh/snippetbox/pkg/models/mysql"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// store configurations for the app
type config struct {
	Addr      string
	StaticDir string
	DSN       string
}

type application struct {
	cfg           *config
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippet       *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// info logger
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	// error logger
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	// configuration
	cfg := new(config)

	// get args
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.DSN, "dsn", "web:pass@tcp(localhost:3306)/snippetbox?parseTime=true", "Mysql database driver DSN (Data Source Name)")

	flag.Parse()

	// setup connection to db
	db, err := openDB(cfg.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	// ensure close is called before exit the program
	defer db.Close()

	// templates
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// start app
	app := application{
		// loggers
		errorLog: errorLog,
		infoLog:  infoLog,

		// db models
		snippet: &mysql.SnippetModel{DB: db},

		// templates
		templateCache: templateCache,

		// config
		cfg: cfg,
	}

	// create a server for custom error logging
	srv := &http.Server{
		Addr:     app.cfg.Addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.Addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
