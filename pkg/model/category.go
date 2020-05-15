package model

import "time"

type Category struct {
	Id        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Position  int       `db:"pos" json:"position"`
	ImageUrl  string    `db:"image_url" json:"image_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
