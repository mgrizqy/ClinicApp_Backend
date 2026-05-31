package models

import "time"

type Appointment struct {
    ID          uint        `gorm:"primaryKey;column:id"`
    PatientID   uint        `gorm:"not null;column:patient_id"`
    Patient     User        `gorm:"foreignKey:PatientID"`
    DoctorID    uint        `gorm:"not null;column:doctor_id"`
    Doctor      User        `gorm:"foreignKey:DoctorID"`
    StartTime   time.Time   `gorm:"not null;column:start_time"`
    EndTime     time.Time   `gorm:"not null;column:end_time"`
}
