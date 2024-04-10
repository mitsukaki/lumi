package models

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

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}