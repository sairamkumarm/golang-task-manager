package main

import (
	"fmt"
	"golang-task-manager/config"
	"golang-task-manager/internal/logger"
	"golang-task-manager/internal/migrations"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	ServerConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:",err)
		os.Exit(1)
	}

	fmt.Println("Config loaded successfully")
	fmt.Println(ServerConfig)
	
	db, err := sqlx.Connect("postgres", ServerConfig.DB_URL)
	if err != nil {
		fmt.Println("failed to connect to db:", err)
		os.Exit(1)
	}
	defer db.Close()

	migrations.RunMigrations(db.DB)

	fmt.Println("Database ready, starting server")

	logger, err := logger.BuildLogger(ServerConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger.Info("Loggered initialised")
}