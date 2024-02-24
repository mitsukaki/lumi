package db

type CouchTable struct {
	baseURL string
	db      *CouchDatabase
}

func (table *CouchTable) Get(docId string) ([]byte, error) {
	return table.db.handle.MakeByteGetRequest(table.baseURL + "/" + docId)
}

func (table *CouchTable) Put(doc interface{}) (*CouchResponse, error) {
	return table.db.handle.MakePostRequest(table.baseURL, doc)
}

func (table *CouchTable) Delete(docId string, rev string) (*CouchResponse, error) {
	return table.db.handle.MakeRequest(DELETE, table.baseURL+"/"+docId+"?rev="+rev, nil)
}
