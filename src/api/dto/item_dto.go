package dto

import (
	"time"

	"github.com/farzadamr/go-backend-api/internal/domain"
	"github.com/farzadamr/go-backend-api/internal/service"
)

// ── Request DTOs (what the client sends) ──────────────────────────────────────

type CreateItemRequest struct {
	Name        string `json:"name"         binding:"required"`
	FileId      int64  `json:"file_id"`
	Description string `json:"description"`
	Price       int    `json:"price"        binding:"required,gt=0"`
	IsAvailable bool   `json:"is_available"`
	Ingredients string `json:"ingredients"`
	CategoryId  int64  `json:"category_id"  binding:"required"`
}

// ToService converts the HTTP DTO → service input model.
// This is the only place json tags touch service layer types.
func (r CreateItemRequest) ToService() service.CreateItemRequest {
	return service.CreateItemRequest{
		Name:        r.Name,
		FileId:      r.FileId,
		Description: r.Description,
		Price:       r.Price,
		IsAvailable: r.IsAvailable,
		Ingredients: r.Ingredients,
		CategoryId:  r.CategoryId,
	}
}

// UpdateItemRequest uses pointers so PATCH works correctly:
// missing fields stay nil and are ignored by the service.
type UpdateItemRequest struct {
	Name        *string `json:"name"`
	FileId      *int64  `json:"file_id"`
	Description *string `json:"description"`
	Price       *int    `json:"price"       binding:"omitempty,gt=0"`
	IsAvailable *bool   `json:"is_available"`
	Ingredients *string `json:"ingredients"`
	CategoryId  *int64  `json:"category_id"`
}

func (r UpdateItemRequest) ToService() service.UpdateItemRequest {
	return service.UpdateItemRequest{
		Name:        r.Name,
		FileId:      r.FileId,
		Description: r.Description,
		Price:       r.Price,
		IsAvailable: r.IsAvailable,
		Ingredients: r.Ingredients,
		CategoryId:  r.CategoryId,
	}
}

type SetAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

// ── Response DTOs (what the client receives) ──────────────────────────────────

type ItemResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	FileID      int64     `json:"file_id"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	IsAvailable bool      `json:"is_available"`
	Ingredients string    `json:"ingredients"`
	CategoryID  int64     `json:"category_id"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

// ItemFromDomain converts domain.Item → the JSON response shape.
func ItemFromDomain(i *domain.Item) ItemResponse {
	return ItemResponse{
		ID:          i.Id,
		Name:        i.Name,
		FileID:      i.FileId,
		Description: i.Description,
		Price:       i.Price,
		IsAvailable: i.IsAvailable,
		Ingredients: i.Ingredients,
		CategoryID:  i.CategoryId,
		CreateAt:    i.CreateAt,
		UpdateAt:    i.UpdateAt,
	}
}

// ── Paginated list response ───────────────────────────────────────────────────

type ItemListResponse struct {
	Items    []ItemResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}
