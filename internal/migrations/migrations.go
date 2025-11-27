package migrations

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatalf("could not create migrate driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "postgres", driver)
    if err != nil {
        log.Fatalf("could not create migrate instance: %v", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("migration failed: %v", err)
    }

    fmt.Println("Migrations applied successfully")
}
