package models

type CouchResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

type UserData struct {
	Username string   `json:"username"`
	Profile  Photo    `json:"profile"`
	Banner   Photo    `json:"banner"`
	Albums   []string `json:"albums"`
}

type DBUser struct {
	ID            string   `json:"_id"`
	Rev           string   `json:"_rev"`
	Email         string   `json:"email"`
	Password      string   `json:"password"`
	Permissions   []string `json:"permissions"`
	PrivateAlbums []string `json:"private_albums"`
	PublicData    UserData `json:"user_data"`

	CouchResponse
}

type Album struct {
	AlbumID      string  `json:"_id"`
	Title        string  `json:"title"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	AuthorUserID string  `json:"author_user_id"`
	CoverPhoto   Photo   `json:"cover_photo"`
	Photos       []Photo `json:"photos"`

	CouchResponse
}

type Photo struct {
	PhotoID string `json:"photo_id"`
	Ratio   int    `json:"ratio"`
}
