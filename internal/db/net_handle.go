package db

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	PUT    = iota
	POST   = iota
	GET    = iota
	DELETE = iota
)

type NetHandle struct {
	authorization string
	httpClient    *http.Client
}

func CreateNetHandle(auth string) *NetHandle {
	return &NetHandle{
		authorization: auth,
		httpClient:    &http.Client{Timeout: time.Second * 15},
	}
}

func (h *NetHandle) MakeGetRequest(url string) (*CouchResponse, error) {
	return h.MakeRequest(GET, url, nil)
}

func (h *NetHandle) MakePostRequest(url string, doc interface{}) (*CouchResponse, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return h.MakeRequest(POST, url, bytes.NewBuffer(jsonData))
}

func (h *NetHandle) MakePutRequest(url string, doc interface{}) (*CouchResponse, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return h.MakeRequest(PUT, url, bytes.NewBuffer(jsonData))
}

func (h *NetHandle) MakeByteGetRequest(url string) ([]byte, error) {
	return h.MakeByteRequest(GET, url, nil)
}

func (h *NetHandle) MakeBytePostRequest(url string, doc interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return h.MakeByteRequest(POST, url, bytes.NewBuffer(jsonData))
}

func (h *NetHandle) MakeBytePutRequest(url string, doc interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	return h.MakeByteRequest(PUT, url, bytes.NewBuffer(jsonData))
}

func (h *NetHandle) MakeRequest(method int, url string, reqBody io.Reader) (*CouchResponse, error) {
	respBody, err := h.MakeByteRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// un-marshal the response body
	var couchResp CouchResponse
	err = json.Unmarshal(respBody, &couchResp)
	if err != nil {
		return nil, err
	}

	return &couchResp, nil
}

func (h *NetHandle) MakeByteRequest(method int, url string, reqBody io.Reader) ([]byte, error) {
	// log the request
	log.Printf("Making request: %s %s", getMethodString(method), url)

	req, err := http.NewRequest(getMethodString(method), url, reqBody)
	if err != nil {
		return nil, err
	}

	// set the headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if h.authorization != "" {
		req.Header.Add("Authorization", "Basic "+h.authorization)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read the response body
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func getMethodString(method int) string {
	if method == POST {
		return "POST"
	}
	if method == GET {
		return "GET"
	}
	if method == DELETE {
		return "DELETE"
	}
	return "PUT"
}
