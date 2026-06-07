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

func RequireRole(allowedRole string) gin.HandlerFunc {

    return func (c *gin.Context) {
        UserRole, exists := c.Get("userRole")
            if !exists {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Authentication credentials missing"})
                return
            }

        userRoleStr, ok := UserRole.(string)

        if !ok {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid credential formatting"})
            return
        }

        if userRoleStr != allowedRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permisions"})
            return
        }

        c.Next()
    }
}

func CORSMiddleware() gin.HandlerFunc{
    return func (c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return 
        }
        c.Next()
    }

}
