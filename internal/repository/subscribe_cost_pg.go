package repository

import (
	"context"
	"effective_mobile/pkg/model"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewSubscribeCostPg(log *zap.Logger, pool *pgxpool.Pool) *SubscribeCostPg {
	return &SubscribeCostPg{log, pool}
}

type SubscribeCostPg struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

func (inst *SubscribeCostPg) SelectTotalCost(ctx context.Context, data *model.SubscriberCostFilter) (float64, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscribes
		WHERE 1=1
	`
	args := []any{}
	argPos := 1

	if data.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, data.UserID)
		argPos++
	}
	if data.StartData != "" {
		query += fmt.Sprintf(" AND start_date >= $%d", argPos)
		args = append(args, data.StartData)
		argPos++
	}
	if data.EndData != "" {
		query += fmt.Sprintf(" AND end_date <= $%d", argPos)
		args = append(args, data.EndData)
		argPos++
	}
	if data.Name != "" {
		query += fmt.Sprintf(" AND name = $%d", argPos)
		args = append(args, data.Name)
		argPos++
	}

	inst.log.Debug("debug query", zap.String("query", query), zap.Any("args", args))

	var total float64
	err := inst.pool.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
