package dto

type CreateAppointmentRequest struct {
	TaskDate    string `json:"task_date" binding:"required,datetime=2006-01-02"`
	Name        string `json:"name" binding:"required,max=80"`
	Description string `json:"description"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`

	AlertType string `json:"alert_type" binding:"required,oneof=none 1_day_before custom_date"`

	AlertDates []string `json:"alert_dates"`
}

type UpdateAppointmentRequest struct {
	TaskDate    string   `json:"task_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Name        string   `json:"name,omitempty" binding:"omitempty,max=80"`
	Description string   `json:"description,omitempty"`
	TimeStart   string   `json:"time_start,omitempty"`
	TimeEnd     string   `json:"time_end,omitempty"`
	AlertType   string   `json:"alert_type,omitempty" binding:"omitempty,oneof=none 1_day_before custom_date"`
	AlertDates  []string `json:"alert_dates,omitempty"`
}