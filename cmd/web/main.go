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
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	mux := http.NewServeMux()

	// configuration
	cfg := new(config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static asse")

	flag.Parse()

	app := application{
		// info logger
		infoLog: log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime),
		// error logger
		errorLog: log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/file", downloadHandler)

	// static files serve
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// create a server for custom error logging
	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: app.errorLog,
		Handler:  mux,
	}

	app.infoLog.Printf("Starting server on %s", cfg.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
