package repository

import (
	"context"
	"effective_mobile/pkg/model"
	"effective_mobile/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewSubscribePg(log *zap.Logger, pool *pgxpool.Pool) *SubscribePg {
	return &SubscribePg{
		log:  log,
		pool: pool,
	}
}

type SubscribePg struct {
	log  *zap.Logger
	pool *pgxpool.Pool
}

func (inst *SubscribePg) SelectOne(ctx context.Context, id string) (*model.Subscribe, error) {
	sub := &model.Subscribe{}

	query := `SELECT id, name, price, user_id, start_date, end_date FROM subscribes WHERE id = $1`

	inst.log.Debug("debug query", zap.String("query", query))

	err := inst.pool.QueryRow(ctx, query, &id).
		Scan(&sub.ID, &sub.Name, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorNotFound
		}
		return nil, err
	}

	return sub, nil
}

func (inst *SubscribePg) SelectMany(ctx context.Context, meta *model.ListData) ([]*model.Subscribe, error) {
	query, args := inst.buildListQuery(
		`SELECT id, name, price, user_id, start_date, end_date FROM subscribes`,
		meta,
	)

	inst.log.Debug("debug query", zap.String("query", query))

	rows, err := inst.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*model.Subscribe
	for rows.Next() {
		sub := &model.Subscribe{}
		if err := rows.Scan(&sub.ID, &sub.Name, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (inst *SubscribePg) Create(ctx context.Context, subscribe *model.Subscribe) error {
	query := `INSERT INTO subscribes (id, name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)`

	inst.log.Debug("debug query", zap.String("query", query))

	_, err := inst.pool.Exec(ctx, query,
		uuid.NewString(),
		subscribe.Name,
		subscribe.Price,
		subscribe.UserID,
		subscribe.StartDate,
		subscribe.EndDate,
	)
	return err
}

func (inst *SubscribePg) Update(ctx context.Context, id string, subscribe *model.Subscribe) error {
	sub := &model.Subscribe{}
	query := `SELECT id, name, price, start_date, end_date FROM subscribes WHERE id = $1`
	err := inst.pool.QueryRow(ctx, query, id).
		Scan(&sub.ID, &sub.Name, &sub.Price, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorNotFound
		}
		return err
	}

	setParts := make([]string, 0)
	args := make([]interface{}, 0)
	argIdx := 1

	if subscribe.Name != "" && subscribe.Name != sub.Name {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, subscribe.Name)
		argIdx++
	}
	if subscribe.Price != 0 && subscribe.Price != sub.Price {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argIdx))
		args = append(args, subscribe.Price)
		argIdx++
	}
	if !time.Time(subscribe.StartDate).IsZero() && subscribe.StartDate != sub.StartDate {
		setParts = append(setParts, fmt.Sprintf("start_date = $%d", argIdx))
		args = append(args, subscribe.StartDate)
		argIdx++
	}
	if !time.Time(subscribe.EndDate).IsZero() && subscribe.EndDate != sub.EndDate {
		setParts = append(setParts, fmt.Sprintf("end_date = $%d", argIdx))
		args = append(args, subscribe.EndDate)
		argIdx++
	}

	if subscribe.UserID != "" && subscribe.UserID != sub.UserID {
		setParts = append(setParts, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, subscribe.UserID)
		argIdx++
	}

	if len(setParts) == 0 {
		return nil
	}

	args = append(args, id)
	query = fmt.Sprintf("UPDATE subscribes SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argIdx)

	inst.log.Debug("debug query", zap.String("query", query))

	cmdTag, err := inst.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return utils.ErrorNotFound
	}
	return nil
}

func (inst *SubscribePg) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM subscribes WHERE id = $1`

	inst.log.Debug("debug query", zap.String("query", query))

	cmdTag, err := inst.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return utils.ErrorNotFound
	}
	return nil
}

func (inst *SubscribePg) buildListQuery(baseQuery string, req *model.ListData) (string, []any) {
	whereParts := []string{}
	args := []any{}
	argIndex := 1

	for _, f := range req.Filters {
		switch strings.ToLower(f.FilterType) {
		case "eq":
			whereParts = append(whereParts, fmt.Sprintf("%s = $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "neq":
			whereParts = append(whereParts, fmt.Sprintf("%s != $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "like":
			whereParts = append(whereParts, fmt.Sprintf("%s ILIKE $%d", f.Field, argIndex))
			args = append(args, "%"+f.Value+"%")
		case "gt":
			whereParts = append(whereParts, fmt.Sprintf("%s > $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "gte":
			whereParts = append(whereParts, fmt.Sprintf("%s >= $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "lt":
			whereParts = append(whereParts, fmt.Sprintf("%s < $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "lte":
			whereParts = append(whereParts, fmt.Sprintf("%s <= $%d", f.Field, argIndex))
			args = append(args, f.Value)
		case "in":
			whereParts = append(whereParts, fmt.Sprintf("%s = ANY($%d)", f.Field, argIndex))
			args = append(args, strings.Split(f.Value, ",")) // "a,b,c" → []string
		}
		argIndex++
	}

	if len(whereParts) > 0 {
		baseQuery += " WHERE " + strings.Join(whereParts, " AND ")
	}

	if len(req.Sorters) > 0 {
		orderParts := []string{}
		for _, s := range req.Sorters {
			dir := "ASC"
			if strings.ToLower(s.Direction) == "desc" {
				dir = "DESC"
			}
			orderParts = append(orderParts, fmt.Sprintf("%s %s", s.Field, dir))
		}
		baseQuery += " ORDER BY " + strings.Join(orderParts, ", ")
	}

	if req.Pagination.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, req.Pagination.Limit, req.Pagination.Offset)
	}

	return baseQuery, args
}
