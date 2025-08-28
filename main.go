package main

import (
    "attendance-payroll/config"
    "attendance-payroll/models"
    "attendance-payroll/routes"
    "github.com/gin-gonic/gin"
    "os"
)

func main() {
    r := gin.Default()

    config.ConnectDB()
    config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})

    routes.SetupRoutes(r)

    // ambil port dari env Railway
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // fallback lokal
    }

    r.Run(":" + port)
}
