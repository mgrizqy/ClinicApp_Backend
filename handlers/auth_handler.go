package handlers

import (
	"backend/db"
	"backend/models"
	"backend/pkg/utils"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

)

type RegisterInput struct {
    Email               string  `json:"email" binding:"required,email"`
    Password            string  `json:"password" binding:"required"`
    Role                string  `json:"role" binding:"required"`
    FirstName           string  `json:"first_name" binding:"required"`
    LastName            string  `json:"last_name" binding:"required"`
    Specialization      string  `json:"specialization"`
    WorkingHoursStart   string  `json:"working_hours_start"`
    WorkingHoursEnd     string  `json:"working_hours_end"`

}

type LoginInput struct {
    Email       string  `json:"email" binding:"required,email"`
    Password    string  `json:"password" binding:"required"`
}





func Register(c *gin.Context) {
    
    var input RegisterInput

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.Role != "patient" && input.Role != "doctor" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalide role specified"})
        return
    }

    if len(input.Password) > 72 || len(input.Password) < 8 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalide password length"})
        return
    }

    if input.Role == "doctor" {
        if input.Specialization == "" || input.WorkingHoursStart == "" || input.WorkingHoursEnd == ""{
            c.JSON(http.StatusBadRequest, gin.H{"error": "Empty Specialization/Working hour start or end"})
            return
        } 

        start, errStart := time.Parse("15:04", input.WorkingHoursStart)
        end, errEnd := time.Parse("15:04", input.WorkingHoursEnd)

        if errStart != nil || errEnd != nil{
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Working hours must match the 'HH:MM' (25 hour) format"})
            return
        }

        if !start.Before(end) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Working hours end must be chronologically after start time"})
            return
        }
    }

    hashedPassword, err := utils.HashPassword(input.Password)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tx := db.DB.Begin()

    defer tx.Rollback()

    newUser := models.User{
        Email: input.Email,
        PasswordHash: hashedPassword,
        Role: input.Role,
    }

    if err := tx.Create(&newUser).Error; err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
        return 
    }

    if input.Role == "patient" {
        newPatient := models.PatientProfile{
            UserID: newUser.ID,
            FirstName: input.FirstName,
            LastName: input.LastName,
        }

        if err := tx.Create(&newPatient).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

    } else if input.Role == "doctor" {
        newDoctor := models.DoctorProfile{
            UserID: newUser.ID,
            FirstName: input.FirstName,
            LastName: input.LastName,
            Specialization: input.Specialization,
            WorkingHoursStart: input.WorkingHoursStart,
            WorkingHoursEnd: input.WorkingHoursEnd,
        }

        if err := tx.Create(&newDoctor).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize to database registration changes"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})

}

func Login(c *gin.Context) {
    var login LoginInput

    if err := c.ShouldBindJSON(&login); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User

    emailCheck := db.DB.Where("email = ?", login.Email).First(&user)

    if emailCheck.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    if !utils.CheckPasswordHash(login.Password, user.PasswordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    token, err := utils.GenerateToken(user.ID, user.Role)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.SetCookie("token", token, 86400, "/", "localhost", false, true) // don't forget to set secure to true

    c.JSON(http.StatusOK, gin.H{"message": "login successful"})
    

}

func Logout(c *gin.Context){
    c.SetCookie("token", "", -1, "/", "localhost", false, true)

    c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
