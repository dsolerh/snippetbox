package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// store configurations for the app
type config struct {
	Addr      string
	StaticDir string
}

type application struct {
	cfg      *config
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// configuration
	app := application{
		// info logger
		infoLog: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		// error logger
		errorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
		cfg:      new(config),
	}

	// get args
	flag.StringVar(&app.cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.cfg.StaticDir, "static-dir", "./ui/static", "Path to static asse")

	flag.Parse()

	// create a server for custom error logging
	srv := &http.Server{
		Addr:     app.cfg.Addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
