package main

import (
    "attendance-payroll/config"
    "attendance-payroll/models"
    "attendance-payroll/routes"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "time"
    "os"
    "fmt"
)

func main() {
    r := gin.Default()

    // ✅ Middleware CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // bisa diganti asal domain, misal "https://yourdomain.com"
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    // ✅ Koneksi database
    config.ConnectDB()
    config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})

    // ✅ Health check
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Attendance Payroll API running 🚀"})
    })

    // ✅ Routes
    routes.SetupRoutes(r)

    // ✅ Jalankan server (baca PORT dari Railway atau default 8080)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("🚀 Server running on port " + port)
    r.Run(":" + port)
}
