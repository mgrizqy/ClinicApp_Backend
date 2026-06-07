package utils

import (
	"errors"
	"time"
)

var(
    ErrPastBooking  = errors.New("cannot book appointment in the past")
    ErrShortNotice  = errors.New("appointments require at least 1 hour of advance notice")
    ErrWeekendBooking  = errors.New("appointments cannot be scheduled on weekends")
    ErrOutsideWorkHours  = errors.New("requested time is outside the doctor's operating hours")
    ErrAppointDuration = errors.New("appointment must be between 1 and 120 minutes")
)

func ValidationBooking(startTime time.Time, durationMinutes int, workStart, workEnd string) error {

    now := time.Now()

    if startTime.Before(now.Add(1 * time.Hour)) {
        return ErrShortNotice
    }

    if durationMinutes < 1 || durationMinutes > 120 {
        return  ErrAppointDuration
    }

    loc, err := time.LoadLocation("America/New_York")

    if err != nil {
        return err
    }

    localTime := startTime.In(loc)

    if localTime.Weekday() == time.Sunday || localTime.Weekday() == time.Saturday  {
        return ErrWeekendBooking
    } 

    localStartStr := localTime.Format("15:04")
    localEndTime := localTime.Add(time.Duration(durationMinutes) * time.Minute)
    localEndStr := localEndTime.Format("15:04")

    if localTime.Day() != localEndTime.Day() {
        return ErrOutsideWorkHours
    }

    if localStartStr < workStart || localEndStr > workEnd {
        return  ErrOutsideWorkHours
    }

    return nil

}
