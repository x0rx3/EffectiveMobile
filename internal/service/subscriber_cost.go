package service

import (
	"context"
	"effective_mobile/internal/repository"
	"effective_mobile/pkg/model"
)

func NewSubscriberCostService(repo repository.SubscriberCostRepository) *SubscriberCost {
	return &SubscriberCost{repo}
}

type SubscriberCost struct {
	repo repository.SubscriberCostRepository
}

func (inst *SubscriberCost) TotalCost(ctx context.Context, data *model.SubscriberCostFilter) (float64, error) {
	return inst.repo.SelectTotalCost(ctx, data)
}
