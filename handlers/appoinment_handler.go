package handlers

import (
	"backend/db"
	"backend/models"
	"backend/pkg/utils"
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

    if err := utils.ValidationBooking(input.StartTime, input.DurationMinutes, doctor.WorkingHoursStart, doctor.WorkingHoursEnd); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func GetAppointments (c *gin.Context) {
    userIDVal, exists := c.Get("userID")

    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
        return
    }

    userID, ok := userIDVal.(uint)
    if !ok {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user identification format"})
        return
    }

   userRole, exists := c.Get("userRole")

   if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication credentials missing"})
        return
   } 

   userRoleStr, ok := userRole.(string)

    if !ok {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid credential formatting"})
        return
    }


    appointments := []models.Appointment{}

        switch userRoleStr {
    case "doctor":
        if err := db.DB.Preload("Patient").Where("doctor_id = ?", userID).Find(&appointments).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return
        }
    case "patient":
        if err := db.DB.Preload("Doctor.DoctorProfile").Where("patient_id = ?", userID).Find(&appointments).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return
        }

    }


    c.JSON(http.StatusOK, appointments)

}
