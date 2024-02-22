package db

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CouchTable struct {
	baseURL string
	db      *CouchDatabase
}

func (table *CouchTable) Put(doc interface{}, docId string) (*CouchResponse, error) {
	// serialize the doc to JSON
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// send a PUT request to the database
	req, err := http.NewRequest(
		"POST",
		table.baseURL+"/"+docId,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := table.db.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// un-marshal the response body
	var couchResp CouchResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &couchResp, nil
}

func (table *CouchTable) Get(docId string) ([]byte, error) {
	// send a GET request to the database
	req, err := http.NewRequest("GET", table.baseURL+"/"+docId, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := table.db.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (table *CouchTable) Delete(docId string, rev string) (*CouchResponse, error) {
	// send a DELETE request to the database
	req, err := http.NewRequest("DELETE", table.baseURL+"/"+docId+"?rev="+rev, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := table.db.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// un-marshal the response body
	var couchResp CouchResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &couchResp, nil
}
