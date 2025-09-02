package transport

import (
	"effective_mobile/internal/service"
	"effective_mobile/internal/transport/handler"
	"net/http"

	"go.uber.org/zap"
)

func NewRouter(log *zap.Logger, serviceManager *service.AppServices) *Router {
	return &Router{
		subscribe:      handler.NewSubscribeHandler(log, serviceManager.SubscriberService),
		subscriberCost: handler.NewSubscriberCost(log, serviceManager.SubscriberCostService),
	}
}

type Router struct {
	subscribe      *handler.Subscribe
	subscriberCost *handler.SubscriberCost
}

func (inst *Router) HandlerMap() map[HandlerMetadata]http.HandlerFunc {
	return map[HandlerMetadata]http.HandlerFunc{
		{"/subscribe", "POST"}:        inst.subscribe.Create,
		{"/subscribe/{id}", "GET"}:    inst.subscribe.Get,
		{"/subscribe", "GET"}:         inst.subscribe.List,
		{"/subscribe/{id}", "DELETE"}: inst.subscribe.Delete,
		{"/subscribe/{id}", "PATCH"}:  inst.subscribe.Update,
		{"/subscribes/cost", "GET"}:   inst.subscriberCost.Get,
	}
}
