package main

import (
	"fmt"
	"log"
	"os"
	"ppob/app/configs"
	"ppob/app/routes"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {

	envFile := ".env"
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(fmt.Sprintf("Error loading :%s file ", envFile));
	}

	appPort := os.Getenv("APP_PORT")

	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBmode := os.Getenv("DB_SSL_MODE")
	DBPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to convert string to int [DB] => %s", err));
	}

	configDB := &configs.ConfigDB{
		DBHost: DBHost,
		DBPort: DBPort,
		DBName: DBName,
		DBUser: DBUser,
		DBPassword: DBPassword,
		DBmode: DBmode,
	};

	db := configs.InitDB(configDB);
	defer db.Close();

	app := routes.InitRoutes(db);
	app.Start(fmt.Sprintf(":%s", appPort));
}