package services

import "github.com/panospet/small-api/pkg/model"

type DbService interface {
	GetProducts(offset int, limit int, orderBy string, asc bool) ([]model.Product, error)
	GetProduct(id string) (model.Product, error)
	AddProduct(product model.Product) (string, error)
	UpdateProduct(product model.Product) error
	DeleteProduct(id string) error
	GetCategories(offset int, limit int, orderBy string, asc bool) ([]model.Category, error)
	GetCategory(id int) (model.Category, error)
	AddCategory(category model.Category) error
	UpdateCategory(category model.Category) error
	DeleteCategory(id int) error
	AddUser(user model.User) error
	UserExists(username string, password string) bool
	AllCategoriesToChan(catC chan model.Category) chan error
	AllProductsToChan(prodC chan model.Product) chan error
}