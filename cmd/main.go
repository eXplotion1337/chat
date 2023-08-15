package main

import (
	"chat/internal/app"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	config, err := app.InitConfig()
	if err != nil {
		log.Fatal("faled to init app config: ", err)
	}

	if err := app.Run(config); err != nil {
		log.Fatal("failed to run app: ", err)
	}
}
