package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/harisoncleytondev/personal-agenda/internal/api/routes/dto"
	"github.com/harisoncleytondev/personal-agenda/internal/model"
	"github.com/harisoncleytondev/personal-agenda/internal/repository"
)

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	emailSvc        *EmailService
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository, emailSvc *EmailService) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		emailSvc:        emailSvc,
	}
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func cleanPtr(s *string) *string {
	if s != nil && *s == "" {
		return nil
	}
	return s
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, userID, userName, userEmail string, req dto.CreateAppointmentRequest) error {
	taskDate, err := time.Parse("2006-01-02", req.TaskDate)
	if err != nil {
		return err
	}

	var alertDates []time.Time
	now := time.Now()
	hoje := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for _, dateStr := range req.AlertDates {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		if parsedDate.Before(hoje) || parsedDate.Equal(hoje) {
			return fmt.Errorf("A data de alerta '%s' não pode ser igual ou anterior a hoje.", dateStr)
		}
		alertDates = append(alertDates, parsedDate)
	}

	newAppointment := model.AppointmentCreate{
		UserID:      userID,
		TaskDate:    taskDate,
		Name:        req.Name,
		Description: nullIfEmpty(req.Description),
		TimeStart:   cleanPtr(req.TimeStart),
		TimeEnd:     cleanPtr(req.TimeEnd),
		AlertType:   req.AlertType,
		AlertDates:  alertDates,
	}

	err = s.appointmentRepo.CreateAppointment(ctx, newAppointment)
	if err != nil {
		return err
	}

	alertInfo := model.AlertInfo{
		TaskName:    req.Name,
		Description: nullIfEmpty(req.Description),
		TimeStart:   cleanPtr(req.TimeStart),
		TaskDate:    taskDate,
		UserEmail:   userEmail,
		UserName:    userName,
	}

	go func() {
		err := s.emailSvc.SendAppointmentAlert(userEmail, userName, alertInfo)
		if err != nil {
			fmt.Printf("ERRO ao enviar e-mail de novo compromisso para %s: %v\n", userEmail, err)
		}
	}()

	return nil
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, id string, userID string, req dto.UpdateAppointmentRequest) error {
	existingApp, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingApp.UserID != userID {
		return errors.New("Acesso negado.")
	}

	taskDate := existingApp.TaskDate
	if req.TaskDate != "" {
		parsed, err := time.Parse("2006-01-02", req.TaskDate)
		if err != nil {
			return err
		}
		taskDate = parsed
	}

	alertDates := existingApp.AlertDates
	if req.AlertDates != nil {
		alertDates = []time.Time{}
		now := time.Now()
		hoje := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		for _, dateStr := range req.AlertDates {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return err
			}

			if parsedDate.Before(hoje) || parsedDate.Equal(hoje) {
				return fmt.Errorf("A data de alerta '%s' não pode ser igual ou anterior a hoje.", dateStr)
			}
			alertDates = append(alertDates, parsedDate)
		}
	}

	name := existingApp.Name
	if req.Name != "" {
		name = req.Name
	}

	alertType := existingApp.AlertType
	if req.AlertType != "" {
		alertType = req.AlertType
	}

	updatedAppointment := model.AppointmentCreate{
		UserID:      userID,
		TaskDate:    taskDate,
		Name:        name,
		Description: nullIfEmpty(req.Description),
		TimeStart:   cleanPtr(req.TimeStart),
		TimeEnd:     cleanPtr(req.TimeEnd),
		AlertType:   alertType,
		AlertDates:  alertDates,
	}

	return s.appointmentRepo.UpdateAppointment(ctx, id, updatedAppointment)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, id string, userID string) error {
	existingApp, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingApp.UserID != userID {
		return errors.New("Acesso negado.")
	}

	return s.appointmentRepo.DeleteAppointment(ctx, id)
}

func (s *AppointmentService) GetAllByUserID(ctx context.Context, userID string) ([]model.Appointment, error) {
	return s.appointmentRepo.GetAllByUserID(ctx, userID)
}

func (s *AppointmentService) GetByID(ctx context.Context, id string, userID string) (*model.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if appointment.UserID != userID {
		return nil, errors.New("Acesso negado.")
	}

	return appointment, nil
}

func (s *AppointmentService) ProcessDailyRoutines(ctx context.Context) error {
	hoje := time.Now()

	alerts, err := s.appointmentRepo.GetAppointmentsToAlert(ctx, hoje)
	if err != nil {
		return err
	}

	userAlerts := make(map[string][]model.AlertInfo)
	userNames := make(map[string]string)

	for _, alert := range alerts {
		userAlerts[alert.UserEmail] = append(userAlerts[alert.UserEmail], alert)
		userNames[alert.UserEmail] = alert.UserName
	}

	for email, alertList := range userAlerts {
		var todayTasks []model.AlertInfo
		var reminderTasks []model.AlertInfo

		for _, alert := range alertList {
			if alert.TaskDate.Format("2006-01-02") == hoje.Format("2006-01-02") {
				todayTasks = append(todayTasks, alert)
			} else {
				reminderTasks = append(reminderTasks, alert)
			}
		}

		err := s.emailSvc.SendDailySummaryAlert(
			email,
			userNames[email],
			todayTasks,
			reminderTasks,
		)

		if err != nil {
			fmt.Printf("ERRO ao enviar e-mail para %s: %v\n", email, err)
		}
	}

	return s.appointmentRepo.DeletePastAppointments(ctx, hoje)
}