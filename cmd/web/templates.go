package main

import (
	"html/template"
	"path/filepath"

	"git.01.alem.school/quazar/forum/pkg/models"
)

type templateData struct {
	Post          models.Post
	Posts         []models.Post
	User          models.User
	Users         []models.User
	Category      models.Category
	Categories    []models.Category
	Comment       models.Comment
	Comments      []models.Comment
	Vote          models.Vote
	Votes         []models.Vote
	IsSession     bool
	CommentNumber string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
