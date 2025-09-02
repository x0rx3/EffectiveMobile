package service

import "effective_mobile/internal/repository"

func NewAppServices(repoManager *repository.PostgresReposManager) *AppServices {
	return &AppServices{
		SubscriberService:     NewSubscribeService(repoManager.SubscriberRepository),
		SubscriberCostService: NewSubscriberCostService(repoManager.SubscriberCostRepository),
	}
}

type AppServices struct {
	SubscriberService     SubscriberService
	SubscriberCostService SubscriberCostService
}
