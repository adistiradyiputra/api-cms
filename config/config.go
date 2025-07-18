package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	AppPort       string
	JWTSecret     string
	RedisURL      string
	RedisPassword string
	RedisDB       int
}

var ENV EnvConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		// log.Printf("Error loading .env file: %v", err)
		log.Fatal("Error loading .env file")
	}

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	ENV = EnvConfig{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		AppPort:       os.Getenv("APP_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisURL:      os.Getenv("REDIS_URL"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
	}

	log.Println("ENV Loaded:")
	log.Printf("Host: %s, Port: %s, User: %s, DB: %s\n", ENV.DBHost, ENV.DBPort, ENV.DBUser, ENV.DBName)
}
