package main

import (
	"backend/db"
	"backend/handlers"
	"backend/middleware"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Println("WARNING: DATABASE_URL not set. Falling back to local development configuration")
        dsn = "host=localhost port=2142 user=postgres password=qwerty123456 dbname=postgres sslmode=disable" // from dbeaver
    }

    var dbConnected bool
    for attempt := 1; attempt <= 5; attempt++ {
        err := db.InitDB(dsn)
        if err == nil {
            dbConnected = true
            log.Println("Successfully connected to the database")
            break
        }
        log.Printf("Database connection attempt %d failed: %v. Retrying in 2 seconds...", attempt, err)
        time.Sleep(2 * time.Second)
    }

    if !dbConnected {
        log.Fatal("Critical: could not connect to database after 5 attempts. Exiting")
    }

    // Open port
    router := gin.Default()

    router.Use(middleware.CORSMiddleware())

    	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Clinic API is running and connected to database",
		})
	})

    router.POST("/api/auth/register", handlers.Register)
    router.POST("/api/auth/login", handlers.Login)
    router.POST("/api/auth/logout", handlers.Logout)
    router.GET("/api/doctors", handlers.GetDoctors)
    router.GET("api/auth/me", middleware.AuthRequire(), handlers.Me)
    router.GET("/api/appointments", middleware.AuthRequire(), handlers.GetAppointments)
    router.POST("/api/appointments", middleware.AuthRequire(), middleware.RequireRole("patient"), handlers.BookAppointment) 
    router.DELETE("/api/users/me", middleware.AuthRequire(), handlers.DeleteMe)
    router.DELETE("/api/users/:id",handlers.DeleteUserByID)


    err := router.Run(":8080")
    if err != nil {
        log.Fatalf("Critical: Server failed to start on port 8080: %v", err)
    }

}
