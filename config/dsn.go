package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Dsn struct {
	User     string
	Password string
	DBName   string
}

func InitDSN() Dsn {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return Dsn{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DBNAME"),
	}
}
