package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	PORT     = 0
	DB_URL   = ""
	DB_DRIVE = ""
	DB_USER  = ""
	DB_PASS  = ""
	DB_HOST  = ""
	DB_PORT  = ""
	DB_NAME  = ""
)

func Load() {
	var err error

	// Load dotEnv
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PORT, err = strconv.Atoi(os.Getenv("BLOG_API_PORT"))
	if err != nil {
		//log.Println(err)
		PORT = 9000
	}

	DB_DRIVE = os.Getenv("BLOG_DB_DRIVE")
	DB_USER = os.Getenv("BLOG_DB_USER")
	DB_PASS = os.Getenv("BLOG_DB_PASSWORD")
	DB_HOST = os.Getenv("BLOG_DB_HOST")
	DB_PORT = os.Getenv("BLOG_DB_PORT")
	DB_NAME = os.Getenv("BLOG_DB_NAME")
	DB_URL = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USER,
		DB_PASS,
		DB_HOST,
		DB_NAME,
	)
}
