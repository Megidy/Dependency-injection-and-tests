package config

import (
	"log"

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
	return Dsn{}
}
