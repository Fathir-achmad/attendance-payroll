package controllers

import (
	"net/http"
	"time"

	"attendance-payroll/config"
	"attendance-payroll/models"

	"github.com/gin-gonic/gin"
)

func GetPayroll(c *gin.Context) {
	userID := c.GetUint("userID")

	// Tentukan periode bulan ini
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, -1)

	// Ambil semua attendance untuk user di bulan ini, preload employee + department
	var attendances []models.Attendance
	if err := config.DB.Preload("Employee").Preload("Employee.Department").
		Where("employee_id = ? AND date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hitung gaji
	dailyRate := 100000.0
	totalSalary := 0.0

	for _, att := range attendances {
		if att.CheckInAt != nil && att.CheckOutAt != nil {
			workHours := att.CheckOutAt.Sub(*att.CheckInAt).Hours()
			if workHours >= 8 {
				totalSalary += dailyRate
			}
		}
	}

	payroll := models.Payroll{
		EmployeeID:  userID,
		Amount:      totalSalary,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	// Isi Employee ke payroll
	if len(attendances) > 0 {
		payroll.Employee = attendances[0].Employee
	}

	c.JSON(http.StatusOK, gin.H{
		"employee_id": userID,
		"payrolls":    []models.Payroll{payroll},
	})
}
