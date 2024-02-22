package models

import "time"

type User struct {
	Username string  `json:"username"`
	Profile  Photo   `json:"profile"`
	Banner   Photo   `json:"banner"`
	Albums   []Album `json:"albums"`
}

type DBUser struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`

	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`

	User
}

type Album struct {
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Description  string    `json:"description"`
	AuthorUserID string    `json:"author_user_id"`
	CoverPhoto   Photo     `json:"cover_photo"`
	Photos       []Photo   `json:"photos"`
}

type Photo struct {
	PhotoID string `json:"photo_id"`
	Ratio   string `json:"ratio"`
}
