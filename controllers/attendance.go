package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "attendance-payroll/config"
    "attendance-payroll/models"
)

func CheckIn(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    employeeID := userID.(uint)
    today := time.Now().Truncate(24 * time.Hour)

    var attendance models.Attendance
    err := config.DB.Where("employee_id = ? AND date = ?", employeeID, today).First(&attendance).Error
    if err == nil && attendance.CheckInAt != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked in"})
        return
    }

    now := time.Now()
    if err != nil {
        attendance = models.Attendance{
            EmployeeID: employeeID,
            Date:       today,
            CheckInAt:  &now,
        }
        config.DB.Create(&attendance)
    } else {
        attendance.CheckInAt = &now
        config.DB.Save(&attendance)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Checked in successfully"})
}

func CheckOut(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    employeeID := userID.(uint)
    today := time.Now().Truncate(24 * time.Hour)

    var attendance models.Attendance
    if err := config.DB.Where("employee_id = ? AND date = ?", employeeID, today).First(&attendance).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Check-in first"})
        return
    }

    if attendance.CheckOutAt != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked out"})
        return
    }

    now := time.Now()
    attendance.CheckOutAt = &now
    config.DB.Save(&attendance)

    c.JSON(http.StatusOK, gin.H{"message": "Checked out successfully"})
}
