package utils

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

var (
	ErrorInvalidAuthData    = errors.New("invalid data")
	ErrorAuthFailed         = errors.New("authorization failed")
	ErrorInvalidToken       = errors.New("invalid token")
	ErrorEmptyID            = errors.New("id cant be empty")
	ErrorNotFound           = errors.New("not found")
	ErrorInvalidAdminToken  = errors.New("invalid admin token")
	ErrorInvalidPassword    = errors.New("invalid password")
	ErrorInvalidLogin       = errors.New("invalid login")
	ErrorInvalidGrant       = errors.New("invalid grant")
	ErrorFilterFormat       = errors.New("invalid filter format")
	ErrorLimitFormat        = errors.New("invalid limit format")
	ErrorEmptyFile          = errors.New("empty file")
	ErrorCacheValue         = errors.New("unxpected type from cache")
	ErrorLoginAlradyExists  = errors.New("user with such a login already has")
	ErrorNoAccess           = errors.New("access denied")
	ErrorInvalidFilterParam = errors.New("invalid filter param")
	ErrorUnxpectedError     = errors.New("Unxpected error")
)

var errorStatusMap = map[error]int{
	ErrorInvalidFilterParam: http.StatusBadRequest,
	ErrorAuthFailed:         http.StatusUnauthorized,
	ErrorEmptyFile:          http.StatusBadRequest,
	ErrorFilterFormat:       http.StatusBadRequest,
	ErrorEmptyID:            http.StatusBadRequest,
	ErrorInvalidGrant:       http.StatusBadRequest,
	ErrorInvalidAuthData:    http.StatusBadRequest,
	ErrorInvalidToken:       http.StatusBadRequest,
	ErrorInvalidAdminToken:  http.StatusBadRequest,
	ErrorInvalidLogin:       http.StatusBadRequest,
	ErrorInvalidPassword:    http.StatusBadRequest,
	ErrorNotFound:           http.StatusNotFound,
	ErrorLoginAlradyExists:  http.StatusConflict,
	ErrorNoAccess:           http.StatusForbidden,
}

func CaseError(w http.ResponseWriter, err error, log *zap.Logger) {
	for target, status := range errorStatusMap {
		if errors.Is(err, target) {
			if log != nil {
				log.Error("error", zap.Int("status", status), zap.Error(target))
			}
			http.Error(w, target.Error(), status)
			return
		}
	}

	if log != nil {
		log.Error("error", zap.Int("status", 500), zap.Error(err))
	}

	http.Error(w, ErrorUnxpectedError.Error(), 500)
}
