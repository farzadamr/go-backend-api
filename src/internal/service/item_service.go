package service

import (
	"context"
	"fmt"

	"github.com/farzadamr/go-backend-api/internal/domain"
	"github.com/farzadamr/go-backend-api/internal/port"
)

type ItemService interface {
	CreateItem(ctx context.Context, req CreateItemRequest) (*domain.Item, error)
	GetItem(ctx context.Context, id int64) (*domain.Item, error)
	UpdateItem(ctx context.Context, id int64, req UpdateItemRequest) (*domain.Item, error)
	DeleteItem(ctx context.Context, id int64) error
	ListItems(ctx context.Context, filter port.ItemFilter) ([]*domain.Item, int64, error)
	SetAvailability(ctx context.Context, id int64, available bool) (*domain.Item, error)
}

type CreateItemRequest struct {
	Name        string
	FileId      int64
	Description string
	Price       int
	IsAvailable bool
	Ingredients string
	CategoryId  int64
}

type UpdateItemRequest struct {
	Name        *string
	FileId      *int64
	Description *string
	Price       *int
	IsAvailable *bool
	Ingredients *string
	CategoryId  *int64
}

type itemService struct {
	itemRepo port.ItemRepository
	// categoryRepo port.CategoryRepository  ← add when you implement Category
	// fileRepo     port.FileRepository      ← add when you implement File
}

func NewItemService(itemRepo port.ItemRepository) ItemService {
	return &itemService{itemRepo: itemRepo}
}

func (s *itemService) CreateItem(ctx context.Context, req CreateItemRequest) (*domain.Item, error) {
	// Business rule: price must be positive
	if req.Price <= 0 {
		return nil, domain.NewInvalidInput("price must be greater than zero")
	}
	if req.Name == "" {
		return nil, domain.NewInvalidInput("name is required")
	}

	item := &domain.Item{
		Name:        req.Name,
		FileId:      req.FileId,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
		Ingredients: req.Ingredients,
		CategoryId:  req.CategoryId,
	}

	created, err := s.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("CreateItem: %w", err)
	}
	return created, nil
}

func (s *itemService) GetItem(ctx context.Context, id int64) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetItem: %w", err)
	}
	return item, nil
}

func (s *itemService) UpdateItem(ctx context.Context, id int64, req UpdateItemRequest) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UpdateItem: %w", err)
	}

	// Only patch fields that were provided (pointer = optional)
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.FileId != nil {
		item.FileId = *req.FileId
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, domain.NewInvalidInput("price must be greater than zero")
		}
		item.Price = *req.Price
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}
	if req.Ingredients != nil {
		item.Ingredients = *req.Ingredients
	}
	if req.CategoryId != nil {
		item.CategoryId = *req.CategoryId
	}

	updated, err := s.itemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("UpdateItem: %w", err)
	}
	return updated, nil
}

func (s *itemService) DeleteItem(ctx context.Context, id int64) error {
	// Check it exists first so we return 404, not a silent no-op
	if _, err := s.itemRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("DeleteItem: %w", err)
	}
	if err := s.itemRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteItem: %w", err)
	}
	return nil
}

func (s *itemService) ListItems(ctx context.Context, filter port.ItemFilter) ([]*domain.Item, int64, error) {
	items, total, err := s.itemRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("ListItems: %w", err)
	}
	return items, total, nil
}

func (s *itemService) SetAvailability(ctx context.Context, id int64, available bool) (*domain.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("SetAvailability: %w", err)
	}
	item.IsAvailable = available
	updated, err := s.itemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("SetAvailability: %w", err)
	}
	return updated, nil
}
