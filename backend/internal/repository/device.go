package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

type DeviceRepo struct {
	pool *pgxpool.Pool
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{pool: pool}
}

func (r *DeviceRepo) Upsert(ctx context.Context, userID uuid.UUID, deviceID, name string) (*models.Device, error) {
	d := &models.Device{UserID: userID, DeviceID: deviceID, Name: name}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO devices (user_id, device_id, name) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, device_id) DO UPDATE SET name = $3, last_seen = NOW()
		 RETURNING id, user_id, device_id, name, last_seen, created_at`,
		userID, deviceID, name,
	).Scan(&d.ID, &d.UserID, &d.DeviceID, &d.Name, &d.LastSeen, &d.CreatedAt)
	return d, err
}

func (r *DeviceRepo) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM devices WHERE user_id = $1`, userID).Scan(&n)
	return n, err
}

func (r *DeviceRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.Device, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, device_id, name, last_seen, created_at FROM devices WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Device
	for rows.Next() {
		var d models.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.DeviceID, &d.Name, &d.LastSeen, &d.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}
