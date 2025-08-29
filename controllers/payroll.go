package controllers

import (
	"attendance-payroll/config"
	"attendance-payroll/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"time"
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

	// parse input bulan
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
	endDate = endDate.AddDate(0, 1, -1) // akhir bulan

	// ambil employee + department
	var employee models.Employee
	if err := config.DB.Preload("Department").First(&employee, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// ambil attendance di rentang waktu
	var attendances []models.Attendance
	if err := config.DB.
		Where("employee_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// gaji harian
	dailyRate := 100000.0
	payrollByMonth := map[string]float64{}

	for _, att := range attendances {
		if att.CheckInAt != nil && att.CheckOutAt != nil {
			workHours := att.CheckOutAt.Sub(*att.CheckInAt).Hours()
			if workHours >= 8 {
				monthKey := att.Date.Format("January 2006")
				payrollByMonth[monthKey] += dailyRate
			}
		}
	}

	// generate list bulan
	var months []time.Time
	for d := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC); d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 1, 0) {
		months = append(months, d)
	}
	sort.Slice(months, func(i, j int) bool { return months[i].Before(months[j]) })

	// buat response payroll per bulan
	var payrolls []gin.H
	for _, m := range months {
		monthKey := m.Format("January 2006")
		amount := payrollByMonth[monthKey]
		payrolls = append(payrolls, gin.H{
			"month":  monthKey,
			"amount": amount,
		})
	}

	// final response
	c.JSON(http.StatusOK, gin.H{
		"employee": gin.H{
			"id":         employee.ID,
			"name":       employee.Name,
			"department": employee.Department.Name,
		},
		"payrolls": payrolls,
	})
}
