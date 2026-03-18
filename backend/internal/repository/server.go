package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"vpn-startup/backend/internal/models"
)

type ServerRepo struct {
	pool *pgxpool.Pool
}

func NewServerRepo(pool *pgxpool.Pool) *ServerRepo {
	return &ServerRepo{pool: pool}
}

func (r *ServerRepo) ListActive(ctx context.Context) ([]models.Server, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, region, host, port, type, reality_pub_key, reality_short_id, reality_sni, is_active, created_at, updated_at
		 FROM servers WHERE is_active = true ORDER BY region, name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Server
	for rows.Next() {
		var s models.Server
		err := rows.Scan(&s.ID, &s.Name, &s.Region, &s.Host, &s.Port, &s.Type, &s.RealityPubKey, &s.RealityShortID, &s.RealitySNI, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *ServerRepo) GetByID(ctx context.Context, id string) (*models.Server, error) {
	s := &models.Server{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, region, host, port, type, reality_pub_key, reality_short_id, reality_sni, is_active, created_at, updated_at
		 FROM servers WHERE id::text = $1`,
		id,
	).Scan(&s.ID, &s.Name, &s.Region, &s.Host, &s.Port, &s.Type, &s.RealityPubKey, &s.RealityShortID, &s.RealitySNI, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}
