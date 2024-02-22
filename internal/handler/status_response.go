package handler

type StatusResponse struct {
	Ok     bool   `json:"ok"`
	Reason string `json:"reason"`
}
