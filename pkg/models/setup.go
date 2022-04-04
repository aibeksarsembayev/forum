package models

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Setup(db *sql.DB) error {
	path := filepath.Join("pkg", "models", "setup.sql")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			return err
		}
	}
	return nil
}
