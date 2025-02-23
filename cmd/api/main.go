package main

import (
	"biblia-be/internal/env"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	app_env := os.Getenv("ENV")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	envFile := ".env.dev"
	if app_env == "dev" {
		envFile = ".env.dev"
	} else if app_env == "prod" {
		envFile = ".env.prod"
	}

	envPath := filepath.Join(dir, envFile)
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		host: env.GetString("HOST", "localhost"),
		db: dbConfig{
			user:         env.GetString("DB_USER", "root"),
			host:         env.GetString("DB_HOST", "127.0.0.1"),
			password:     env.GetString("DB_PASSWORD", "znoksy139"),
			db_name:      env.GetString("DB_NAME", "biblia_db"),
			db_addr:      env.GetString("DB_ADDR", "3306"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetInt("DB_MAX_IDLE_TIME", 15),
		},
	}

	app := &application{
		config: cfg,
	}

	app.run()
}
