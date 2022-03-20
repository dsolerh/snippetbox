package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/file", app.downloadHandler)

	// static files serve
	fileServer := http.FileServer(http.Dir(app.cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
