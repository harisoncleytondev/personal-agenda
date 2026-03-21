package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB * pgxpool.Pool

func Connect() {
	var err error

	DB, err = pgxpool.New(context.Background(), getDatabaseLink())

	if err != nil {
        log.Fatal("Erro ao conectar no banco:", err)
    }

    err = DB.Ping(context.Background())
    if err != nil {
        log.Fatal("Banco não respondeu:", err)
    }

	log.Println("Conectado ao banco de dados.")
}