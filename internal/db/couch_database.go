package db

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

type CouchDBConfig struct {
	UseAuth  bool
	UseHTTPS bool
	Host     string
	Port     int
	Username string
	Password string
}

type CouchDatabase struct {
	baseURL string
	handle  *NetHandle
}

func CreateDatabase(config *CouchDBConfig) *CouchDatabase {
	scheme := "http://"
	if config.UseHTTPS {
		scheme = "https://"
	}

	url := scheme + config.Host + ":" + strconv.Itoa(config.Port) + "/"
	var auth string
	if config.UseAuth {
		fullToken := config.Username + ":" + config.Password
		auth = base64.URLEncoding.EncodeToString([]byte(fullToken))
		// print the auth string
		fmt.Println(auth)
	} else {
		auth = ""
	}

	return &CouchDatabase{
		baseURL: url,
		handle:  CreateNetHandle(auth),
	}
}

func (db *CouchDatabase) CreateTable(name string) (*CouchTable, error) {
	url := db.baseURL + name
	_, err := db.handle.MakeRequest(PUT, url, nil)
	if err != nil {
		return nil, err
	}

	table := &CouchTable{
		baseURL: url,
		db:      db,
	}

	return table, nil
}
