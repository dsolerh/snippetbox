package main

import "dsolerh.projects/snippetbox/pkg/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
