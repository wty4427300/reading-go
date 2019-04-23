package main

import (
	"io"
	"log"
	"os"

	"github.com/developer-learning/reading-go/blob/master/examples/ginexamples/pkg/config"
	"github.com/developer-learning/reading-go/blob/master/examples/ginexamples/pkg/http"
	"github.com/developer-learning/reading-go/blob/master/examples/ginexamples/pkg/postgres"
	"github.com/developer-learning/reading-go/blob/master/examples/ginexamples/pkg/service/userservice"
)

func main() {
	c := config.GetConfig()
	postgresConfig := postgres.DBConfig{
		Host:     c.PGHost,
		Port:     c.PGPort,
		User:     c.PGUser,
		Password: c.PGPassword,
		Name:     c.PGDBName,
	}

	repository := postgres.Initialize(postgresConfig)
	repository.AutoMigrate()

	var logDst io.Writer
	if c.LogFile == "" {
		logDst = os.Stdout
	} else {
		file, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("error opening logfile %s: %v", c.LogFile, err)
		}
		defer file.Close()
		logDst = file
	}

	server := http.AppServer{
		Logger:      log.New(logDst, "", log.LstdFlags),
		UserService: userservice.New(repository.UserRepository),
	}
	server.Run()
}
