package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/panospet/small-api/internal/config"
	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
)

func main() {
	var workers int
	var amount int

	flag.IntVar(&workers, "workers", 10, "number of workers")
	flag.IntVar(&amount, "amount", 1000, "total amount of products to add")
	flag.Parse()


	conf := config.NewConfig()
	db, err := services.NewDb(conf.MysqlPath)
	if err != nil {
		panic(err)
	}
	redis, err := cache.NewRedisCache(conf.RedisPath)
	populateDb(db, workers, amount)
	populateRedis(db, redis, workers)
}

func populateRedis(db *services.AppDb, cacher *cache.RedisCacher, workers int) {
	prodC := make(chan model.Product)
	catC := make(chan model.Category)
	wg := sync.WaitGroup{}

	startRedis := time.Now()

	// fill channels with products and categories
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = db.AllProductsToChan(prodC)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = db.AllCategoriesToChan(catC)
	}()

	// populate redis with categories
	wg.Add(1)
	go func() {
		defer wg.Done()
		for cat := range catC {
			catBytes, _ := json.Marshal(cat)
			_ = cacher.SetCategory(fmt.Sprintf("%d", cat.Id), string(catBytes))
		}
	}()

	// populate redis with products
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for prod := range prodC {
				prodBytes, _ := json.Marshal(prod)
				_ = cacher.SetProduct(prod.Id, string(prodBytes))
			}
		}()
	}

	wg.Wait()

	redisDuration := time.Since(startRedis)
	fmt.Println("Redis populated. Time:", redisDuration.String())
}

func populateDb(db *services.AppDb, workers int, amount int) {
	rand.Seed(time.Now().UnixNano())
	possibleCategories := []string{"sports", "house", "garden", "electronics", "games", "food", "drinks", "furniture",
		"space", "mobile", "movies", "tv", "pc", "books", "groceries", "devices", "music", "instruments"}
	startDb := time.Now()
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
	products := make(chan model.Product)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pr := range products {
				_, err := db.AddProduct(pr)
				if err != nil {
					log.Println("error adding product to database:", err)
				}
			}
		}()
	}
	for i := 0; i < amount; i++ {
		products <- model.Product{
			CategoryId:  rand.Intn(len(possibleCategories)) + 1,
			Title:       fmt.Sprintf("product %d", i),
			ImageUrl:    fmt.Sprintf("http://www.bestprice.gr/product%d.png", i),
			Price:       float32(rand.Intn(200)) + rand.Float32(),
			Description: fmt.Sprintf("Description for product %d", i),
		}
	}
	close(products)
	wg.Wait()

	dbDuration := time.Since(startDb)
	fmt.Println("Database populated. Time:", dbDuration.String())
}
