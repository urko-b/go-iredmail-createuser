package main

import (
	"fmt"
	"log"

	"iredmail-create-email-account/controllers"
	"iredmail-create-email-account/middleware"
	"iredmail-create-email-account/pkg/remote_ssh"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var cfg struct {
	RemoteUser     string   `required:"true" split_words:"true"`
	RemoteServer   string   `required:"true" split_words:"true"`
	RemoteKeyPath  string   `required:"true" split_words:"true"`
	RemotePort     string   `required:"true" split_words:"true"`
	ApiPort        string   `required:"true" split_words:"true"`
	SupportEmail   string   `required:"true" split_words:"true"`
	TrustedProxies []string `required:"true" split_words:"true"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(
			fmt.Errorf("environment variable ENV is empty and an error occurred while loading the .env file"),
		)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	config := remote_ssh.Config{
		User:    cfg.RemoteUser,
		Server:  cfg.RemoteServer,
		KeyPath: cfg.RemoteKeyPath,
		Port:    cfg.RemotePort,
	}

	router := gin.Default()
	router.SetTrustedProxies(cfg.TrustedProxies)
	router.Use(middleware.Json())
	router.Use(middleware.Error(cfg.SupportEmail))

	user_controller := controllers.NewUser(config)
	router.POST("/user", user_controller.CreateUser)

	err = router.Run(fmt.Sprintf("localhost:%s", cfg.ApiPort))
	if err != nil {
		panic(err)
	}
}
