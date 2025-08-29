package routes

import (
    "attendance-payroll/controllers"
    "attendance-payroll/middlewares"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        // Auth
        api.POST("/auth/register", controllers.Register)
        api.POST("/auth/login", controllers.Login)

        // Attendance (JWT)
        auth := api.Group("/attendances")
        auth.Use(middlewares.AuthMiddleware())
        {
            auth.POST("/checkin", controllers.CheckIn)
            auth.POST("/checkout", controllers.CheckOut)
            auth.POST("/payroll", controllers.GetPayroll)
        }

        // CRUD Profile
        user := api.Group("/employee")
        user.Use(middlewares.AuthMiddleware())
        {
            user.GET("/me", controllers.GetProfile)
            user.PATCH("/me", controllers.UpdateProfile)
            user.DELETE("/me", controllers.DeleteAccount)
        }


        // Department (Basic Auth)
        dept := api.Group("/departments")
        dept.Use(middlewares.AuthMiddleware())
        {
            dept.POST("", controllers.CreateDepartment)
            dept.GET("", controllers.GetDepartments)
            dept.PATCH("/:id", controllers.UpdateDepartment)
            dept.DELETE("/:id", controllers.DeleteDepartment)
        }
    }
}
