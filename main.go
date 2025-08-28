package main

import (
	"attendance-payroll/config"
	"attendance-payroll/models"
	"attendance-payroll/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.ConnectDB()
	config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})

	// seed default dept biar register bisa jalan
	config.DB.FirstOrCreate(&models.Department{}, models.Department{Name: "Information Technology"})

	routes.SetupRoutes(r)

	// ambil port dari env Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback lokal
	}

	r.Run(":" + port)
}
