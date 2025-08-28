package main

import (
    "attendance-payroll/config"
    "attendance-payroll/models"
    "attendance-payroll/routes"
    "github.com/gin-gonic/gin"
    "os"
    "fmt"
)

func main() {
    r := gin.Default()

    config.ConnectDB()
    config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})

    // health check
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Attendance Payroll API running 🚀"})
    })

    routes.SetupRoutes(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("🚀 Server running on port " + port)
    r.Run(":" + port)
}
