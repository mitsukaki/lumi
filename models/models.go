package models

type CouchResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type Album struct {
	AlbumID      string  `json:"_id"`
	AuthorUserID string  `json:"author_user_id"`
	CoverPhoto   Photo   `json:"cover_photo"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	Title        string  `json:"title"`
	Photos       []Photo `json:"photos"`

	CouchResponse
}

type DBAlbum struct {
	AlbumID      string  `json:"_id"`
	Rev          string  `json:"_rev"`
	AuthorUserID string  `json:"author_user_id"`
	CoverPhoto   Photo   `json:"cover_photo"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	Title        string  `json:"title"`
	Photos       []Photo `json:"photos"`

	CouchResponse
}

type AlbumPutRequest struct {
	CoverPhoto  Photo  `json:"cover_photo"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Private     bool   `json:"private"`
}

type Photo struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PhotoID     string  `json:"photo_id"`
	Ratio       float64 `json:"ratio"`
	Row         int     `json:"row"`
}
