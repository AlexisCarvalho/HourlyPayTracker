package config

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dbPath := GetEnv("DB_PATH", "calculaPagamento.db")
	fmt.Printf("Using database path: %q\n", dbPath)

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Falha ao conectar ao banco de dados: %v", err))
	}

	DB = database
}
