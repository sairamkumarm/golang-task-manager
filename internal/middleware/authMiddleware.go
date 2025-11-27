package middleware

import (
    "golang-task-manager/internal/utils"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(auth, "Bearer ")

        userID, err := utils.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        c.Set("user_id", userID)
        c.Next()
    }
}
