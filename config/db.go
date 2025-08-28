package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    _ "github.com/lib/pq"
)

var DB *gorm.DB

func ConnectDB() {
    // ambil dari env
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=fathiras1905 dbname=absensi port=5432 sslmode=disable"
    }

    // extract info dasar
    host := "localhost"
    user := "postgres"
    pass := "fathiras1905"
    port := "5432"
    dbName := "absensi"

    // connect dulu ke default db postgres
    defaultDSN := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=postgres sslmode=disable",
        host, user, pass, port)
    defaultDB, err := sql.Open("postgres", defaultDSN)
    if err != nil {
        log.Fatal("Failed to connect to default database:", err)
    }
    defer defaultDB.Close()

    // cek apakah database absensi ada
    var exists bool
    err = defaultDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName).Scan(&exists)
    if err != nil {
        log.Fatal("Failed to check database existence:", err)
    }

    // kalau belum ada → buat
    if !exists {
        _, err = defaultDB.Exec("CREATE DATABASE " + dbName)
        if err != nil {
            log.Fatal("Failed to create database:", err)
        }
        log.Printf("✅ Database %s created successfully", dbName)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db
    fmt.Println("✅ Connected to database:", dbName)
}
