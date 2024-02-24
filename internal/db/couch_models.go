package db

type CouchResponse struct {
	Ok     bool   `json:"ok"`
	ID     string `json:"id"`
	Rev    string `json:"rev"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type CouchDatabaseRequest struct {
	Database string `json:"db"`
}	