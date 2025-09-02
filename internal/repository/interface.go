package repository

import (
	"context"
	"effective_mobile/pkg/model"
)

type SubscriberRepository interface {
	SelectOne(ctx context.Context, id string) (*model.Subscribe, error)
	SelectMany(ctx context.Context, meta *model.ListData) ([]*model.Subscribe, error)
	Create(ctx context.Context, subscribe *model.Subscribe) error
	Update(ctx context.Context, id string, subscribe *model.Subscribe) error
	Delete(ctx context.Context, id string) error
}

type SubscriberCostRepository interface {
	SelectTotalCost(ctx context.Context, data *model.SubscriberCostFilter) (float64, error)
}
