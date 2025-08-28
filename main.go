package main

import (
	"attendance-payroll/config"
	"attendance-payroll/models"
	"attendance-payroll/routes"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	DB  *sql.DB
	err error
)

func main() {
	r := gin.Default()

	config.ConnectDB()
	config.DB.AutoMigrate(&models.Department{}, &models.Employee{}, &models.Attendance{}, &models.Payroll{})
	// Build connection string
	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	// Open DB
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	routes.SetupRoutes(r)

	r.Run(":8080")
	r.Run(":" + os.Getenv("PORT"))

}
