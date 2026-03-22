package repository

import (
	"context"

	"github.com/harisoncleytondev/personal-agenda/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepository struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}


func (r *AppointmentRepository) CreateAppointment(ctx context.Context, ap model.AppointmentCreate) error {
	query := `
		INSERT INTO appointments (user_id, task_date, name, description, time_start, time_end, alert_type, alert_dates)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query, ap.UserID, ap.TaskDate,ap.Name, ap.Description, ap.TimeStart, ap.TimeEnd, ap.AlertType, ap.AlertDates)

	if err != nil {
        return err
    }

    return nil
}

func (r *AppointmentRepository) GetAllByUserID(ctx context.Context, userID string) ([]model.Appointment, error) {
	query := `
		SELECT id, user_id, task_date, name, description, time_start, time_end, alert_type, alert_dates, created_at, updated_at 
		FROM appointments 
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []model.Appointment
	for rows.Next() {
		var ap model.Appointment
		err := rows.Scan(
			&ap.ID, &ap.UserID, &ap.TaskDate, &ap.Name, &ap.Description,
			&ap.TimeStart, &ap.TimeEnd, &ap.AlertType, &ap.AlertDates,
			&ap.CreatedAt, &ap.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, ap)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id string) (*model.Appointment, error) {
	query := `
		SELECT id, user_id, task_date, name, description, time_start, time_end, alert_type, alert_dates, created_at, updated_at 
		FROM appointments 
		WHERE id = $1
	`

	var ap model.Appointment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ap.ID, &ap.UserID, &ap.TaskDate, &ap.Name, &ap.Description,
		&ap.TimeStart, &ap.TimeEnd, &ap.AlertType, &ap.AlertDates,
		&ap.CreatedAt, &ap.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ap, nil
}

func (r *AppointmentRepository) UpdateAppointment(ctx context.Context, id string, ap model.AppointmentCreate) error {
	query := `
		UPDATE appointments 
		SET task_date = $1, name = $2, description = $3, time_start = $4, time_end = $5, alert_type = $6, alert_dates = $7, updated_at = NOW()
		WHERE id = $8
	`

	_, err := r.db.Exec(ctx, query, ap.TaskDate, ap.Name, ap.Description, ap.TimeStart, ap.TimeEnd, ap.AlertType, ap.AlertDates, id)
	return err
}

func (r *AppointmentRepository) DeleteAppointment(ctx context.Context, id string) error {
	query := `
		DELETE FROM appointments 
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}