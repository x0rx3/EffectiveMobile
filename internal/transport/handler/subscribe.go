package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"effective_mobile/internal/service"
	"effective_mobile/internal/transport/dto"
	"effective_mobile/pkg/model"
	"effective_mobile/pkg/utils"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewSubscribeHandler(log *zap.Logger, service service.SubscriberService) *Subscribe {
	return &Subscribe{log, service}
}

type Subscribe struct {
	log     *zap.Logger
	service service.SubscriberService
}

// Create godoc
// @Summary Create new subscribe
// @Description Create new subscribe
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param data body dto.CreateUpdateSubscribeDto true "Subscribe data"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /subscribe [post]
func (inst *Subscribe) Create(w http.ResponseWriter, req *http.Request) {
	sub, err := inst.readModelFromReq(req)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	if err := inst.service.Create(context.Background(), sub); err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// Get godoc
// @Summary Get subscribe by ID
// @Description Get subscribe by ID
// @Tags Subscribe
// @Produce json
// @Param id path string true "Subscribe ID"
// @Success 200 {object} dto.ReadSubscribeDto
// @Failure 400
// @Failure 404
// @Router /subscribe/{id} [get]
func (inst *Subscribe) Get(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		utils.CaseError(w, utils.ErrorEmptyID, inst.log)
		return
	}

	sub, err := inst.service.Get(context.Background(), id)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&dto.ReadSubscribeDto{
		ID:        sub.ID,
		Name:      sub.Name,
		UserID:    sub.UserID,
		Price:     sub.Price,
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
	})
}

// List godoc
// @Summary Get list of subscribes
// @Description Get list of subscribes with pagination, filtering and sorting
// @Tags Subscribe
// @Produce json
// @Access json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param sort query string false "Sort fields, e.g. name:asc,price:desc"
// @Param filters query string false "Filters in the form [{"Field":"user_id","Value":"123","FilterType":"eq"}]"
// @Success 200 {array} dto.ReadSubscribeDto
// @Failure 400
// @Failure 500
// @Router /subscribe [get]
func (inst *Subscribe) List(w http.ResponseWriter, r *http.Request) {
	meta, err := inst.parseListRequestFromQuery(r)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	subs, err := inst.service.List(context.Background(), meta)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	dtoSub := make([]dto.ReadSubscribeDto, len(subs))
	for i, sub := range subs {
		dtoSub[i] = dto.ReadSubscribeDto{
			ID:        sub.ID,
			Name:      sub.Name,
			UserID:    sub.UserID,
			Price:     sub.Price,
			StartDate: sub.StartDate,
			EndDate:   sub.EndDate,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoSub)
}

// Update godoc
// @Summary Update subscribe by ID
// @Description Update subscribe by ID
// @Tags Subscribe
// @Accept json
// @Produce json
// @Param id path string true "Subscribe ID"
// @Param data body dto.CreateUpdateSubscribeDto true "Subscribe data"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /subscribe/{id} [patch]
func (inst *Subscribe) Update(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		utils.CaseError(w, utils.ErrorEmptyID, inst.log)
		return
	}

	sub, err := inst.readModelFromReq(req)
	if err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	if err := inst.service.Update(context.Background(), id, sub); err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// Delete godoc
// @Summary Delete subscribe by ID
// @Description Delete subscribe by ID
// @Tags Subscribe
// @Produce json
// @Param id path string true "Subscribe ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /subscribe/{id} [delete]
func (inst *Subscribe) Delete(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		utils.CaseError(w, utils.ErrorEmptyID, inst.log)
		return
	}

	if err := inst.service.Delete(context.Background(), id); err != nil {
		utils.CaseError(w, err, inst.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (inst *Subscribe) readModelFromReq(req *http.Request) (*model.Subscribe, error) {
	sub := &dto.CreateUpdateSubscribeDto{}
	rawData, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(rawData, sub); err != nil {
		return nil, err
	}

	return &model.Subscribe{
		Name:      sub.Name,
		UserID:    sub.UserID,
		Price:     sub.Price,
		StartDate: sub.StartDate,
		EndDate:   sub.EndDate,
	}, nil
}

func (inst *Subscribe) parseListRequestFromQuery(r *http.Request) (*model.ListData, error) {
	req := &model.ListData{}
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit == 0 {
		limit = 10
	}
	req.Pagination = model.Pagination{Limit: limit, Offset: offset}

	if filtersJSON := q.Get("filters"); filtersJSON != "" {
		var filters []model.Filter
		if err := json.Unmarshal([]byte(filtersJSON), &filters); err != nil {
			return nil, err
		}
		req.Filters = filters
	}

	if sortParam := q.Get("sort"); sortParam != "" {
		for _, s := range strings.Split(sortParam, ",") {
			parts := strings.SplitN(s, ":", 2)
			field := parts[0]
			dir := "asc"
			if len(parts) == 2 {
				dir = parts[1]
			}
			req.Sorters = append(req.Sorters, model.Sorter{
				Field:     field,
				Direction: dir,
			})
		}
	}

	return req, nil
}
