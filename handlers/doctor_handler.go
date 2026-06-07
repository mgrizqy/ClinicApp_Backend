package handlers

import (
	"backend/db"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDoctors(c *gin.Context) {
    doctors := []models.DoctorProfile{}

    if err := db.DB.Preload("User").Find(&doctors).Error; err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusOK, doctors)
}
