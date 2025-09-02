package service

import (
	"context"
	"effective_mobile/internal/repository"
	"effective_mobile/pkg/model"
	"effective_mobile/pkg/utils"
)

func NewSubscribeService(
	repo repository.SubscriberRepository,
) *Subscribe {
	return &Subscribe{repo}
}

type Subscribe struct {
	repo repository.SubscriberRepository
}

func (inst *Subscribe) Get(ctx context.Context, id string) (*model.Subscribe, error) {
	return inst.repo.SelectOne(ctx, id)
}

func (inst *Subscribe) List(ctx context.Context, meta *model.ListData) ([]*model.Subscribe, error) {
	filterField := map[string]struct{}{
		"name":       {},
		"start_date": {},
		"end_date":   {},
		"user_id":    {},
		"price":      {},
	}

	for _, f := range meta.Filters {
		if _, ok := filterField[f.Field]; !ok {
			return nil, utils.ErrorInvalidFilterParam
		}
	}
	return inst.repo.SelectMany(ctx, meta)

}
func (inst *Subscribe) Create(ctx context.Context, subscribe *model.Subscribe) error {
	return inst.repo.Create(ctx, subscribe)
}

func (inst *Subscribe) Update(ctx context.Context, id string, subscribe *model.Subscribe) error {
	return inst.repo.Update(ctx, id, subscribe)
}

func (inst *Subscribe) Delete(ctx context.Context, id string) error {
	return inst.repo.Delete(ctx, id)
}
