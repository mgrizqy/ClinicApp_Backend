package handlers

import (
	"backend/db"
	"backend/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookAppointmentInput struct {
    DoctorID        uint        `json:"doctor_id" binding:"required"`
    StartTime       time.Time   `json:"start_time"`
    DurationMinutes int         `json:"duration_minutes"`
}

func BookAppointment (c *gin.Context) {
    var input BookAppointmentInput

    if err:= c.ShouldBindJSON(&input); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userIDVal, exists := c.Get("userID")

    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(uint)

    if !ok {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    endTime := input.StartTime.Add(time.Duration(input.DurationMinutes) * time.Minute)

    var doctor models.DoctorProfile

    if err := db.DB.Where("user_id = ?", input.DoctorID).First(&doctor).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Doctor profile not found"})
            return
        }

        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return

    }

    startTimeStr := input.StartTime.Format("15:04")
    endTimeStr := endTime.Format("15:04")

    if startTimeStr < doctor.WorkingHoursStart || endTimeStr > doctor.WorkingHoursEnd {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "appointment is outside working hours"})
        return
    }

    var conflictCount int64

    if err := db.DB.Model(&models.Appointment{}).Where(
        "doctor_id = ? AND start_time < ? AND end_time > ?",
        input.DoctorID,
        endTime,
        input.StartTime,
         ).Count(&conflictCount).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return
         }

    if conflictCount > 0 {
        c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "This appointment time slot is already booked"})
        return
    }

    newAppointment := models.Appointment{
        PatientID: userID,
        DoctorID: input.DoctorID,
        StartTime: input.StartTime,
        EndTime: endTime,
    }

    if err := db.DB.Create(&newAppointment).Error; err != nil{
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Successfully book an appointment"})

}
