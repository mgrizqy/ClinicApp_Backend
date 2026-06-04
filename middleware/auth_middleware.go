package middleware

import (
	"backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)


func AuthRequire() gin.HandlerFunc{
    return func(c *gin.Context) {
        tokenString, err := c.Cookie("token")
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication token cookie missing or expired"})
            return
        }

        verifiedToken, err := utils.VerifyToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired authentication token"})
            return
        }

        c.Set("userID", verifiedToken.UserID)
        c.Set("userRole", verifiedToken.Role)

        c.Next()
    }
}
