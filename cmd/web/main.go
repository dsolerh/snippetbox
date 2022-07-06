package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// my package for snippet related functionalities
	"dsolerh/snippetbox/pkg/models"
	"dsolerh/snippetbox/pkg/models/mysql"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// store configurations for the app
type config struct {
	Addr      string
	StaticDir string
	DSN       string
	Secret    string
}

type application struct {
	cfg           *config
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      models.ISnippetModel
	users         models.IUserModel
	templateCache map[string]*template.Template
}

type contextKey string

var contextKeyUser = contextKey("user")

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
	flag.StringVar(&cfg.Secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")

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

	session := sessions.New([]byte(cfg.Secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// start app
	app := application{
		// loggers
		errorLog: errorLog,
		infoLog:  infoLog,

		// db models
		snippets: &mysql.SnippetModel{DB: db},
		users:    &mysql.UserModel{DB: db},

		// templates
		templateCache: templateCache,

		// session
		session: session,

		// config
		cfg: cfg,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// create a server for custom error logging
	srv := &http.Server{
		Addr: app.cfg.Addr,
		// logging
		ErrorLog: app.errorLog,
		// routes
		Handler: app.routes(),
		// tls config
		TLSConfig: tlsConfig,
		// timeouts
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
