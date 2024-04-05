package api

type StatusResponse struct {
	Ok     bool   `json:"ok"`
	Reason string `json:"reason"`
}

type AlbumPutResponse struct {
	Ok bool   `json:"ok"`
	ID string `json:"album_id"`
}
