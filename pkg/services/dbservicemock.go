package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/panospet/small-api/pkg/model"
	"math/rand"
	"time"
)

type DbServiceMock struct {
	Products   []model.Product
	Categories []model.Category
}

func NewMockDb() *DbServiceMock {
	categories := generateCategories()
	products := generateProducts()
	return &DbServiceMock{
		Products:   products,
		Categories: categories,
	}
}

func (s *DbServiceMock) GetProducts(offset int, limit int, orderBy string, asc bool) ([]model.Product, error) {
	return s.Products, nil
}

func (s *DbServiceMock) GetProduct(id string) (model.Product, error) {
	return s.Products[0], nil
}

func (s *DbServiceMock) AddProduct(product model.Product) (string, error) {
	id := uuid.New().String()
	product.Id = id
	s.Products = append(s.Products, product)
	return id, nil
}

func (s *DbServiceMock) UpdateProduct(product model.Product) error {
	s.Products[0] = product
	return nil
}

func (s *DbServiceMock) DeleteProduct(id string) error {
	index := 0
	for i, p := range s.Products {
		if p.Id == id {
			index = i
			break
		}
	}
	if index == 0 {
		return errors.New("product not found")
	}
	removeProdFromSlice(s.Products, index)
	return nil
}

func (s *DbServiceMock) GetCategories(offset int, limit int, orderBy string, asc bool) ([]model.Category, error) {
	return s.Categories, nil
}

func (s *DbServiceMock) GetCategory(id int) (model.Category, error) {
	for _, c := range s.Categories {
		if c.Id == id {
			return c, nil
		}
	}
	return model.Category{}, errors.New("product not found")
}

func (s *DbServiceMock) AddCategory(category model.Category) error {
	category.Id = s.Categories[len(s.Categories)-1].Id + 1
	s.Categories = append(s.Categories, category)
	return nil
}

func (s *DbServiceMock) UpdateCategory(category model.Category) error {
	for i, c := range s.Categories {
		if c.Id == category.Id {
			s.Categories[i] = category
			return nil
		}
	}
	return errors.New("category not found")
}

func (s *DbServiceMock) DeleteCategory(id int) error {
	index := 0
	for i, c := range s.Categories {
		if c.Id == id {
			index = i
			break
		}
	}
	if index == 0 {
		return errors.New("category not found")
	}
	removeCatFromSlice(s.Categories, index)
	return nil
}

func (s *DbServiceMock) AddUser(user model.User) error {
	return errors.New("method not implemented")
}

func (s *DbServiceMock) UserExists(username string, password string) bool {
	return false
}

func (s *DbServiceMock) AllCategoriesToChan(catC chan model.Category) chan error {
	errC := make(chan error)
	return errC
}

func (s *DbServiceMock) AllProductsToChan(prodC chan model.Product) chan error {
	errC := make(chan error)
	return errC
}

var possibleCategories = []string{"sports", "house", "garden", "electronics", "games", "food", "drinks", "furniture",
	"space", "mobile", "movies", "tv", "pc", "books", "groceries", "devices", "music", "instruments"}

func generateProducts() []model.Product {
	var products []model.Product
	for i := 0; i < 200; i++ {
		products = append(products, generateProduct(i))
	}
	return products
}

func generateCategories() []model.Category {
	var categories []model.Category
	for i := 0; i < len(possibleCategories); i++ {
		categories = append(categories, model.Category{
			Id:        i,
			Title:     possibleCategories[i],
			Position:  i,
			ImageUrl:  fmt.Sprintf("http://www.bestprice.gr/cat%d.png", i),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
	return categories
}

func generateProduct(i int) model.Product {
	id := uuid.New().String()
	rand.Seed(time.Now().UnixNano())
	return model.Product{
		Id:          id,
		CategoryId:  rand.Intn(len(possibleCategories)) + 1,
		Title:       fmt.Sprintf("product%d", i),
		ImageUrl:    fmt.Sprintf("http://www.bestprice.gr/product%d.png", i),
		Price:       float32(rand.Intn(200)) + rand.Float32(),
		Description: fmt.Sprintf("Description product %d", i),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Category:    model.Category{},
	}
}

func removeProdFromSlice(slice []model.Product, s int) []model.Product {
	return append(slice[:s], slice[s+1:]...)
}

func removeCatFromSlice(slice []model.Category, s int) []model.Category {
	return append(slice[:s], slice[s+1:]...)
}
