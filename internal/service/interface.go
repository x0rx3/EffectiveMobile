package service

import (
	"context"
	"effective_mobile/pkg/model"
)

type SubscriberService interface {
	Get(ctx context.Context, id string) (*model.Subscribe, error)
	List(ctx context.Context, meta *model.ListData) ([]*model.Subscribe, error)
	Create(ctx context.Context, subscribe *model.Subscribe) error
	Update(ctx context.Context, id string, subscribe *model.Subscribe) error
	Delete(ctx context.Context, id string) error
}

type SubscriberCostService interface {
	TotalCost(ctx context.Context, data *model.SubscriberCostFilter) (float64, error)
}
