package config

import (
    "fmt"
    "log"
    "os"
    "strings"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("❌ DATABASE_URL is not set in environment")
    }

    // Railway biasanya kasih prefix "postgres://", ganti ke "postgresql://"
    if strings.HasPrefix(dsn, "postgres://") {
        dsn = strings.Replace(dsn, "postgres://", "postgresql://", 1)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect to database: %v", err)
    }

    DB = db
    fmt.Println("✅ Connected to database")
}
