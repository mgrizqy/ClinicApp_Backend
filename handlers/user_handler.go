package handlers

import (
	"backend/db"
	"backend/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DoctorDetailsResponse struct {
	Specialization    string `json:"specialization"`
	WorkingHoursStart string `json:"working_hours_start"`
	WorkingHoursEnd   string `json:"working_hours_end"`
}

type UnifiedMeResponse struct {
	UserID        uint                   `json:"user_id"`
	Email         string                 `json:"email"`
	Role          string                 `json:"role"`
	FirstName     string                 `json:"first_name"`
	LastName      string                 `json:"last_name"`
	DoctorDetails *DoctorDetailsResponse `json:"doctor_details"` // pointer so it can be 'null' for patients!
}

func Me (c *gin.Context) {
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
    
    var user models.User

    if err := db.DB.First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User session no longer valid"})
            return
        }

        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
        return
    }

    var patient models.PatientProfile
    var doctor models.DoctorProfile
    var response UnifiedMeResponse

    switch user.Role {
        case "patient":
            if err := db.DB.Where("user_id = ?", user.ID).First(&patient).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User id doesn't exist"})
                    return
                }

                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
                return
            }

            response.FirstName = patient.FirstName
            response.LastName = patient.LastName

        case "doctor":
            if err := db.DB.Where("user_id = ?", user.ID).First(&doctor).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User id doesn't exist"})
                    return
                }

                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
                return
            }

            response.FirstName = doctor.FirstName
            response.LastName = doctor.LastName
            response.DoctorDetails = &DoctorDetailsResponse{
                Specialization: doctor.Specialization,
                WorkingHoursStart: doctor.WorkingHoursStart,
                WorkingHoursEnd: doctor.WorkingHoursEnd,
            }
    }

    response.UserID = user.ID
    response.Email = user.Email
    response.Role = user.Role

    c.JSON(http.StatusOK, response)
    

}

func DeleteMe (c *gin.Context) {
    userIDVal, exists := c.Get("userID")
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    userID, ok := userIDVal.(uint)
    
    if !ok {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }


    var user models.User

   if err := db.DB.First(&user, userID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User session no longer valid"})
            return
        }

        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
        return
    }

    tx := db.DB.Begin()

    defer tx.Rollback()

    switch user.Role {
        case "doctor" :
            if err := tx.Where("user_id = ?", userID).Delete(&models.DoctorProfile{}).Error; err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
        case "patient" :
            if err := tx.Where("user_id = ?", userID).Delete(&models.PatientProfile{}).Error; err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
    }

    if err := tx.Where("patient_id = ? OR doctor_id = ?", userID, userID).Delete(&models.Appointment{}).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return 
    }

    if err := tx.Delete(&user).Error; err != nil{
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
    }

    if err:= tx.Commit().Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
    }

    c.SetCookie("token", "", -1, "/", "localhost", false, true)

    c.JSON(http.StatusOK, gin.H{"message": "Account has been deleted successfully"})

}

func DeleteUserByID (c *gin.Context) {
    targetIDStr := c.Param("id")
    id, err := strconv.Atoi(targetIDStr)

    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Your request was malformed, send a valid numeric ID"})
        return
    }

    var user models.User

    if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error occurred"})
        return
    }

    tx := db.DB.Begin()

    defer tx.Rollback()

    switch user.Role {
        case "doctor" :
            if err := tx.Where("user_id = ?", id).Delete(&models.DoctorProfile{}).Error; err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
        case "patient" :
            if err := tx.Where("user_id = ?", id).Delete(&models.PatientProfile{}).Error; err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
    }

    if err := tx.Where("patient_id = ? OR doctor_id = ?", id, id).Delete(&models.Appointment{}).Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return 
    }

    if err := tx.Delete(&user).Error; err != nil{
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
    }

    if err:= tx.Commit().Error; err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
    }


    c.JSON(http.StatusOK, gin.H{"message": "Account has been deleted successfully"})

}
