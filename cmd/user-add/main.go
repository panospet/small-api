package main

import (
	"flag"
	"log"

	"github.com/panospet/small-api/internal/config"
	"github.com/panospet/small-api/pkg/model"
	"github.com/panospet/small-api/pkg/services"
)

func main() {
	var username string
	var password string

	flag.StringVar(&username, "username", "", "admin username")
	flag.StringVar(&password, "password", "", "admin password")
	flag.Parse()

	if username == "" || password == "" {
		log.Fatalln("please give username and password")
	}

	conf := config.NewConfig()
	db, err := services.NewDb(conf.MysqlPath)
	if err != nil {
		panic(err)
	}

	user := model.User{
		Username: username,
		Password: password,
	}

	err = db.AddUser(user)
	if err != nil {
		panic(err)
	}
}
