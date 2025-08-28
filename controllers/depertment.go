package controllers

import (
    "attendance-payroll/config"
    "attendance-payroll/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

// Create Department
func CreateDepartment(c *gin.Context) {
    var input struct {
        Name string `json:"name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    department := models.Department{Name: input.Name}
    if err := config.DB.Create(&department).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, department)
}

// Get All Departments
func GetDepartments(c *gin.Context) {
    var departments []models.Department
    if err := config.DB.Preload("Employees").Find(&departments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, departments)
}

// Update Department
func UpdateDepartment(c *gin.Context) {
    id := c.Param("id")
    var department models.Department
    if err := config.DB.First(&department, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
        return
    }

    var input struct {
        Name string `json:"name" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    department.Name = input.Name
    config.DB.Save(&department)
    c.JSON(http.StatusOK, department)
}

// Delete Department
func DeleteDepartment(c *gin.Context) {
    id := c.Param("id")
    if err := config.DB.Delete(&models.Department{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Department deleted"})
}
