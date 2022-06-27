package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// create a middleware chain
	standardMiddleware := alice.New(app.panicRecover, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	// routes
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	mux.Get("/file", http.HandlerFunc(app.downloadHandler))

	// static files serve
	fileServer := http.FileServer(http.Dir(app.cfg.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
