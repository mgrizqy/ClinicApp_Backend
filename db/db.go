package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is our shared database pool pointer.
var DB *gorm.DB

// InitDB attempts to connect to PostgreSQL.
func InitDB(dsn string) error {
	var err error
	
	// Line 1: We attempt to open the database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Line 2: If an error occurred, we return it immediately.
	if err != nil {
		return err
	}

	// Extract the generic database interface *sql.DB from GORM.
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Set defensive limits on our connection pool.
	sqlDB.SetMaxOpenConns(10) // Limit total simultaneous connections to 10
	sqlDB.SetMaxIdleConns(5)  // Keep up to 5 idle connections open in the pool

	// Line 3: If we reach this point, everything succeeded! We return nil.
	return nil
}
