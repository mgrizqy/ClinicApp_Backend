package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID              uint        `gorm:"primaryKey;column:id"`
    Email           string      `gorm:"not null;column:email"`
    PasswordHash    string      `gorm:"not null;column:password_hash" json:"-"`
    Role            string      `gorm:"not null;column:role"` // strictly "patient" or "doctor"
    CreatedAt       time.Time   `gorm:"column:created_at"`
    UpdatedAt       time.Time   `gorm:"column:updated_at"`
    DeletedAt       gorm.DeletedAt   `gorm:"column:deleted_at;index"`
    DoctorProfile   *DoctorProfile  `gorm:"foreignKey:UserID"`
    PatientProfile  *PatientProfile `gorm:"foreignKey:UserID"`
}

type PatientProfile struct {
    ID          uint        `gorm:"primaryKey;column:id"`
    UserID      uint        `gorm:"not null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    User        User        `gorm:"foreignKey:UserID"`
    FirstName   string      `gorm:"not null;column:first_name"`
    LastName    string      `gorm:"not null;column:last_name"`
    CreatedAt   time.Time   `gorm:"column:created_at"`
    UpdatedAt   time.Time   `gorm:"column:updated_at"`
}

type DoctorProfile struct {
    ID                  uint        `gorm:"primaryKey;column:id"`
    UserID              uint        `gorm:"not null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    User                User        `gorm:"foreignKey:UserID"`
    FirstName           string      `gorm:"not null;column:first_name"`
    LastName            string      `gorm:"not null;column:last_name"`
    Specialization      string      `gorm:"not null;column:specialization"`
    WorkingHoursStart   string      `gorm:"not null;column:working_hours_start"`
    WorkingHoursEnd     string      `gorm:"not null;column:working_hours_end"`
    CreatedAt           time.Time   `gorm:"column:created_at"`
    UpdatedAt           time.Time   `gorm:"column:updated_at"`
}
