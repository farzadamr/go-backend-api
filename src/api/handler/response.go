package handler

import (
	"errors"
	"net/http"

	"github.com/farzadamr/go-backend-api/internal/domain"
	"github.com/gin-gonic/gin"
)

// ── Standard response envelopes ───────────────────────────────────────────────

type errorBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type successBody struct {
	Data any `json:"data"`
}

type listBody struct {
	Data     any   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// ── Respond helpers (called from every handler) ───────────────────────────────

func respondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, successBody{Data: data})
}

func respondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, successBody{Data: data})
}

func paginatedRespondList(c *gin.Context, data any, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, listBody{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func respondList(c *gin.Context, data any, total int64) {
	c.JSON(http.StatusOK, listBody{
		Data:  data,
		Total: total,
	})
}

func respondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// respondError is the ONLY place that translates domain errors → HTTP codes.
// Every handler calls this; nothing else touches net/http status codes.
func respondError(c *gin.Context, err error) {
	code, body := mapError(err)
	c.JSON(code, body)
}

func respondBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, errorBody{
		Error:   "bad_request",
		Message: msg,
	})
}

// mapError: domain sentinel → HTTP status code.
// Uses errors.Is so wrapped errors (fmt.Errorf("...: %w", err)) still match.
func mapError(err error) (int, errorBody) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, errorBody{
			Error:   "not_found",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, errorBody{
			Error:   "conflict",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest, errorBody{
			Error:   "invalid_input",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrBusinessRule):
		return http.StatusUnprocessableEntity, errorBody{
			Error:   "business_rule_violation",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, errorBody{
			Error:   "unauthorized",
			Message: err.Error(),
		}
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, errorBody{
			Error:   "forbidden",
			Message: err.Error(),
		}
	default:
		// Never leak internal errors to the client
		return http.StatusInternalServerError, errorBody{
			Error:   "internal_error",
			Message: "an unexpected error occurred",
		}
	}
}
