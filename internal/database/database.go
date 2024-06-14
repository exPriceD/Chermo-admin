package database

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

var DB *sqlx.DB

func InitDB() (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	log.Println("Connection string:", connStr)

	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}
