package model

import "time"

type Appointment struct {
	ID          string      `json:"id" db:"id"`
	UserID      string      `json:"user_id" db:"user_id"`
	TaskDate    time.Time   `json:"task_date" db:"task_date"`
	Name        string      `json:"name" db:"name"`
	Description *string     `json:"description" db:"description"`
	TimeStart   *string     `json:"time_start" db:"time_start"` 
	TimeEnd     *string     `json:"time_end" db:"time_end"`
	AlertType   string      `json:"alert_type" db:"alert_type"`
	AlertDates  []time.Time `json:"alert_dates" db:"alert_dates"` 
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

type AppointmentCreate struct {
	ID          string      `db:"id"`
	UserID      string      `db:"user_id"`
	TaskDate    time.Time   `db:"task_date"`
	Name        string      `db:"name"`
	Description *string     `db:"description"`
	TimeStart   *string     `db:"time_start"`
	TimeEnd     *string     `db:"time_end"`
	AlertType   string      `db:"alert_type"`
	AlertDates  []time.Time `db:"alert_dates"`
}

type AlertInfo struct {
	TaskName    string
	Description *string
	TimeStart   *string
	TaskDate    time.Time
	UserEmail   string
	UserName    string
}