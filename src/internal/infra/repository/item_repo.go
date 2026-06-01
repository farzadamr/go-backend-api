package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/farzadamr/go-backend-api/internal/domain"
	"github.com/farzadamr/go-backend-api/internal/infra/database"
	"github.com/farzadamr/go-backend-api/internal/port"
)

// ── DB model (GORM tags live HERE, never on domain.Item) ─────────────────────

type itemModel struct {
	Id          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	FileId      int64  `gorm:"column:file_id"`
	Description string
	Price       int  `gorm:"not null"`
	IsAvailable bool `gorm:"default:true"`
	Ingredients string
	CategoryId  int64 `gorm:"column:category_id;not null"`
	CreateAt    time.Time
	UpdateAt    time.Time
}

func (itemModel) TableName() string { return "items" }

// ── Constructor ───────────────────────────────────────────────────────────────

type itemRepository struct {
	db       *gorm.DB
	preloads []database.PreloadEntity
}

// NewItemRepository returns port.ItemRepository — not *itemRepository.
// The caller (app.go) never sees the concrete type.
func NewItemRepository(preloads []database.PreloadEntity) port.ItemRepository {
	return &itemRepository{
		db:       database.GetDb(),
		preloads: preloads,
	}
}

// ── Interface implementation ──────────────────────────────────────────────────

func (r *itemRepository) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	model := toItemModel(item)
	model.CreateAt = time.Now().UTC()
	model.UpdateAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	return toDomainItem(&model), nil
}

func (r *itemRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	var model itemModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Translate DB error → typed domain error
		return nil, domain.NewNotFound("item not found")
	}
	if err != nil {
		return nil, err
	}
	return toDomainItem(&model), nil
}

func (r *itemRepository) Update(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	model := toItemModel(item)
	model.UpdateAt = time.Now().UTC()
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	return toDomainItem(&model), nil
}

func (r *itemRepository) Delete(ctx context.Context, id int64) error {
	// GORM soft-delete: sets DeletedAt, keeps the row
	return r.db.WithContext(ctx).Delete(&itemModel{}, id).Error
}

func (r *itemRepository) List(ctx context.Context, filter port.ItemFilter) ([]*domain.Item, int64, error) {
	query := r.db.WithContext(ctx).Model(&itemModel{})

	// Apply optional filters
	if filter.CategoryId != 0 {
		query = query.Where("category_id = ?", filter.CategoryId)
	}
	if filter.IsAvailable != nil {
		query = query.Where("is_available = ?", filter.IsAvailable)
	}

	// Count before pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginatio

	var models []itemModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*domain.Item, len(models))
	for i := range models {
		items[i] = toDomainItem(&models[i])
	}
	return items, total, nil
}

// ── Mapping helpers (the only place DB ↔ domain conversion happens) ───────────

func toItemModel(i *domain.Item) itemModel {
	return itemModel{
		Id:          i.Id,
		Name:        i.Name,
		FileId:      i.FileId,
		Description: i.Description,
		Price:       i.Price,
		IsAvailable: i.IsAvailable,
		Ingredients: i.Ingredients,
		CategoryId:  i.CategoryId,
	}
}

func toDomainItem(m *itemModel) *domain.Item {
	return &domain.Item{
		Id:          m.Id,
		Name:        m.Name,
		FileId:      m.FileId,
		Description: m.Description,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		Ingredients: m.Ingredients,
		CategoryId:  m.CategoryId,
		CreateAt:    m.CreateAt,
		UpdateAt:    m.UpdateAt,
	}
}
