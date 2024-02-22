package db

import (
	"net/http"
)

type CouchDatabase struct {
	baseURL    string
	httpClient *http.Client
}

func CreateDatabase(baseURL string) *CouchDatabase {
	return &CouchDatabase{
		baseURL:    baseURL,
		httpClient: &http.Client{
			Timeout: 10,
		},
	}
}

func (db *CouchDatabase) CreateTable(name string) *CouchTable {
	return &CouchTable{
		baseURL: db.baseURL + "/" + name,
		db:      db,
	}
}
