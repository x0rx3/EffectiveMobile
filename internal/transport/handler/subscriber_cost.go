package handler

import (
	"effective_mobile/internal/service"
	"effective_mobile/internal/transport/dto"
	"effective_mobile/pkg/model"
	"effective_mobile/pkg/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func NewSubscriberCost(log *zap.Logger, service service.SubscriberCostService) *SubscriberCost {
	return &SubscriberCost{log, service}
}

type SubscriberCost struct {
	log     *zap.Logger
	service service.SubscriberCostService
}

// Get godoc
// @Summary      Get total subscription cost
// @Description  Returns the total subscription cost for a user with optional filters by subscription name and date range
// @Tags         subscriber-cost
// @Accept       json
// @Produce      json
// @Param        user_id     query   string  false   "User ID"
// @Param        name        query   string  false  "Subscription name"
// @Param        start_date  query   string  false  "Start date (YYYY-MM-DD)"
// @Param        end_date    query   string  false  "End date (YYYY-MM-DD)"
// @Success      200  {object}  dto.SubscribeCostDto
// @Failure      400
// @Failure      500
// @Router       /subscribes/cost [get]
func (inst *SubscriberCost) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	data := &model.SubscriberCostFilter{
		UserID:    q.Get("user_id"),
		Name:      q.Get("name"),
		StartData: q.Get("start_date"),
		EndData:   q.Get("end_date"),
	}

	total, err := inst.service.TotalCost(r.Context(), data)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&dto.SubscribeCostDto{Total: total})
}
