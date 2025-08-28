package middlewares

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func GetJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "mysecret"
    }
    return []byte(secret)
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return GetJWTSecret(), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // ✅ Ambil claims dari token
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            if id, ok := claims["id"].(float64); ok {
                c.Set("userID", uint(id)) // simpan userID ke context
            }
        }

        c.Next()
    }
}
