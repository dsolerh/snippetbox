package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// create a middleware chain
	standardMiddleware := alice.New(app.panicRecover, app.logRequest, secureHeaders)

	mux := pat.New()

	// routes
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	mux.Get("/file", http.HandlerFunc(app.downloadHandler))

	// static files serve
	fileServer := http.FileServer(http.Dir(app.cfg.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
