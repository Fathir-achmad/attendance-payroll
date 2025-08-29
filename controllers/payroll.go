package controllers

import (
	"attendance-payroll/config"
	"attendance-payroll/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"sort"
)

type PayrollRequest struct {
	Start string `json:"start" binding:"required"` // format MM-YYYY
	End   string `json:"end" binding:"required"`   // format MM-YYYY
}

func GetPayroll(c *gin.Context) {
	userID := c.GetUint("userID")

	var input PayrollRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// parse MM-YYYY → time.Time (awal bulan)
	startDate, err := time.Parse("01-2006", input.Start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start format (use MM-YYYY)"})
		return
	}

	endDate, err := time.Parse("01-2006", input.End)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end format (use MM-YYYY)"})
		return
	}
	// akhir bulan endDate
	endDate = endDate.AddDate(0, 1, -1)

	var attendances []models.Attendance
	if err := config.DB.Preload("Employee").
		Where("employee_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dailyRate := 100000.0
	payrollByMonth := map[string]float64{}

	for _, att := range attendances {
		if att.CheckInAt != nil && att.CheckOutAt != nil {
			workHours := att.CheckOutAt.Sub(*att.CheckInAt).Hours()
			if workHours >= 8 {
				monthKey := att.Date.Format("January 2006") // contoh: "April 2025"
				payrollByMonth[monthKey] += dailyRate
			}
		}
	}

	// urutkan berdasarkan waktu
	months := make([]string, 0, len(payrollByMonth))
	for m := range payrollByMonth {
		months = append(months, m)
	}
	sort.Slice(months, func(i, j int) bool {
		ti, _ := time.Parse("January 2006", months[i])
		tj, _ := time.Parse("January 2006", months[j])
		return ti.Before(tj)
	})

	var payrolls []gin.H
	for _, m := range months {
		payrolls = append(payrolls, gin.H{
			"month":  m,
			"amount": payrollByMonth[m],
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"employee_id": userID,
		"payrolls":    payrolls,
	})
}
