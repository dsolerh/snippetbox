package main

import (
	"flag"
	"log"
	"net/http"
)

// store configurations for the app
type Config struct {
	Addr      string
	StaticDir string
}

func main() {
	mux := http.NewServeMux()

	// configuration
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static asse")

	flag.Parse()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/file", downloadHandler)

	// static files serve
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, mux)
	log.Fatal(err)
}
