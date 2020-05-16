package main

import (
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

	user := model.User{
		Username:  "admin", // todo username and password using args
		Password:  "admin",
	}

	err = db.AddUser(user)
	if err != nil {
		panic(err)
	}
}
