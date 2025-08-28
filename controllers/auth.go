package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    "attendance-payroll/config"
    "attendance-payroll/models"
    "attendance-payroll/middlewares"
)

// CREATE (Register)
func Register(c *gin.Context) {
    var input struct {
        Name         string `json:"name"`
        Username     string `json:"username"`
        Password     string `json:"password"`
        DepartmentID uint   `json:"department_id"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

    employee := models.Employee{
        Name:         input.Name,
        Username:     input.Username,
        PasswordHash: string(hashed),
        DepartmentID: input.DepartmentID,
    }

    if err := config.DB.Create(&employee).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LOGIN
func Login(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    var employee models.Employee
    if err := config.DB.Where("username = ?", input.Username).First(&employee).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    if bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(input.Password)) != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":       employee.ID,
        "username": employee.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, _ := token.SignedString(middlewares.GetJWTSecret())

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// READ (Get Profile)
func GetProfile(c *gin.Context) {
    userID := c.GetUint("userID")
    var employee models.Employee
    if err := config.DB.Preload("Department").First(&employee, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, employee)
}

// UPDATE (Update Profile)
func UpdateProfile(c *gin.Context) {
    userID := c.GetUint("userID")
    var employee models.Employee
    if err := config.DB.First(&employee, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    var input struct {
        Name         string `json:"name"`
        Password     string `json:"password"`
        DepartmentID uint   `json:"department_id"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if input.Name != "" {
        employee.Name = input.Name
    }
    if input.Password != "" {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
        employee.PasswordHash = string(hashed)
    }
    if input.DepartmentID != 0 {
        employee.DepartmentID = input.DepartmentID
    }

    config.DB.Save(&employee)
    c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// DELETE (Delete Account)
func DeleteAccount(c *gin.Context) {
    userID := c.GetUint("userID")
    if err := config.DB.Delete(&models.Employee{}, userID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete account"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
