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

	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, -1)

	// 1. Cek apakah payroll sudah ada di DB
	var payroll models.Payroll
	if err := config.DB.Preload("Employee").Preload("Employee.Department").
		Where("employee_id = ? AND period_start = ? AND period_end = ?", userID, periodStart, periodEnd).
		First(&payroll).Error; err == nil {
		// Payroll sudah ada → langsung return
		c.JSON(http.StatusOK, payroll)
		return
	}

	// 2. Kalau belum ada, hitung dari attendance
	var attendances []models.Attendance
	if err := config.DB.Preload("Employee").Preload("Employee.Department").
		Where("employee_id = ? AND date BETWEEN ? AND ?", userID, periodStart, periodEnd).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dailyRate := 100000
	totalSalary := 0

	for _, att := range attendances {
		if att.CheckInAt != nil && att.CheckOutAt != nil {
			workHours := att.CheckOutAt.Sub(*att.CheckInAt).Hours()
			if workHours >= 8 {
				totalSalary += dailyRate
			}
		}
	}

	payroll = models.Payroll{
		EmployeeID:  userID,
		Amount:      float64(totalSalary),
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	if len(attendances) > 0 {
		payroll.Employee = attendances[0].Employee
	}

	// 3. Simpan payroll baru ke DB
	if err := config.DB.Create(&payroll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payroll"})
		return
	}

	c.JSON(http.StatusOK, payroll)
}
