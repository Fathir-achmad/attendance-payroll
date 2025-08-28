package main

import (
    "attendance-payroll/config"
    "attendance-payroll/models"
    "attendance-payroll/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    config.ConnectDB()
    config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})

    routes.SetupRoutes(r)

    r.Run(":8080")
}
