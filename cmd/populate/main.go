package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/panospet/small-api/internal/config"
	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
)

func main() {
	conf := config.NewConfig()
	db, err := services.NewDb(conf.MysqlPath)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	possibleCategories := []string{"sports", "house", "garden", "electronics", "games", "food", "drinks"}
	for i := range possibleCategories {
		err := db.AddCategory(model.Category{
			Title:    possibleCategories[i],
			Position: rand.Intn(20) + 1,
			ImageUrl: fmt.Sprintf("http://www.bestprice.gr/%s.png", possibleCategories[i]),
		})
		if err != nil {
			panic(fmt.Sprintf("Error while creating category: %s", err.Error()))
		}
	}
	for i := 0; i < 100; i++ {
		_, err := db.AddProduct(model.Product{
			CategoryId:  rand.Intn(len(possibleCategories)) + 1,
			Title:       fmt.Sprintf("product %d", i),
			ImageUrl:    fmt.Sprintf("http://www.bestprice.gr/product%d.png", i),
			Price:       float32(rand.Intn(200)) + rand.Float32(),
			Description: fmt.Sprintf("Description for product %d", i),
		})
		if err != nil {
			panic(fmt.Sprintf("Error while creating product: %s", err.Error()))
		}
	}
}
