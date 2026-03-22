package main

import (
	"log"

	"github.com/harisoncleytondev/personal-agenda/config"
	"github.com/harisoncleytondev/personal-agenda/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Falhou ao carregar as variaveis de ambiente.")
	}

	config.Connect()

	r := api.SetupRouter()

	if err := r.Run(":" + config.GetPort()); err != nil {
		log.Fatalln("Falhou ao iniciar o servidor.")
	} 

}