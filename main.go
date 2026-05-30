package main

import(
    "log"
    "os"
    "time"
    "github.com/gin-gonic/gin"
    "backend/db"
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

    router.GET("/ping", func(c *gin.Context){
        c.JSON(200, gin.H{
            "status": "healthy",
            "message": "Clinic API is running and connected to database", 
        })
    })

    err := router.Run(":8080")
    if err != nil {
        log.Fatalf("Critical: Server failed to start on port 8080: %v", err)
    }

}
