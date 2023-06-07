package configs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ConfigDB struct {
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBmode     string
}

func InitDB(config *ConfigDB) *sql.DB {
	dataSource := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`, config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBmode);
	db, err := sql.Open("postgres", dataSource);
	if err != nil {
		log.Fatal(fmt.Sprintf("DB Error Open => %s", err));
	}

	db.SetMaxIdleConns(5);
	db.SetMaxOpenConns(10);
	db.SetConnMaxIdleTime(5 * time.Minute);
	db.SetConnMaxLifetime(20 * time.Minute);

	if err = db.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("Connection DB Error => %s", err));
	}

	return db;
}