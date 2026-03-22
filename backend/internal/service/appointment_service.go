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
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{appointmentRepo: appointmentRepo}
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

func (s *AppointmentService) CreateAppointment(ctx context.Context, userID string, req dto.CreateAppointmentRequest) error {
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

	err = s.appointmentRepo.UpdateAppointment(ctx, id, updatedAppointment)
	if err != nil {
		return err
	}

	return nil
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, id string, userID string) error {
	existingApp, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingApp.UserID != userID {
		return errors.New("Acesso negado.")
	}

	err = s.appointmentRepo.DeleteAppointment(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AppointmentService) GetAllByUserID(ctx context.Context, userID string) ([]model.Appointment, error) {
	appointments, err := s.appointmentRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return appointments, nil
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