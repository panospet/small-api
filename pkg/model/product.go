package model

import "time"

type Product struct {
	Id          string    `db:"id" json:"id"`
	CategoryId  int       `db:"category_id" json:"category_id"`
	Title       string    `db:"title" json:"title"`
	ImageUrl    string    `db:"image_url" json:"image_url"`
	Price       float32   `db:"price" json:"price"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
