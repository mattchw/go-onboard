package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongo_auth := ""

	if os.Getenv("MONGO_USER") != "" && os.Getenv("MONGO_PASSWORD") != "" {
		mongo_auth = os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@"
	}

	// return formatted mongo uri (mongo_protocol://mongo_user:mongo_password@mongo_host:mongo_port/mongo_database)
	return fmt.Sprintf("%s://%s%s:%s/%s", os.Getenv("MONGO_PROTOCOL"), mongo_auth, os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"), os.Getenv("MONGO_DATABASE"))
}

func EnvMongoDatabase() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_DATABASE")
}
