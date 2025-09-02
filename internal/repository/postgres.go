package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPostgresReposManager(log *zap.Logger) *PostgresReposManager {
	return &PostgresReposManager{log: log}
}

type PostgresReposManager struct {
	log                      *zap.Logger
	pool                     *pgxpool.Pool
	SubscriberRepository     SubscriberRepository
	SubscriberCostRepository SubscriberCostRepository
}

func (inst *PostgresReposManager) Connect(dsn string) (err error) {
	inst.pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return err
	}

	inst.SubscriberRepository = NewSubscribePg(inst.log, inst.pool)
	inst.SubscriberCostRepository = NewSubscribeCostPg(inst.log, inst.pool)

	return nil
}

func (inst *PostgresReposManager) Close() {
	inst.pool.Close()
}
