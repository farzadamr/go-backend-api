package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/farzadamr/go-backend-api/api/dto"
	"github.com/farzadamr/go-backend-api/internal/dependency"
	"github.com/farzadamr/go-backend-api/internal/port"
	"github.com/farzadamr/go-backend-api/internal/service"
)

type ItemHandler struct {
	itemSvc service.ItemService // interface — never the concrete *itemService
}

func NewItemHandler() *ItemHandler {
	itemRepo := dependency.GetItemRepository()
	return &ItemHandler{itemSvc: service.NewItemService(itemRepo)}
}

// ─────────────────────────────────────────────────────────────────────────────
// Every handler follows the same 4-step pattern:
//   1. Parse   — bind JSON / path params / query params
//   2. Validate — ShouldBindJSON covers binding tags; add extra checks if needed
//   3. Call service — pass service-layer input types (not DTOs)
//   4. Respond — convert domain result → DTO, write JSON
// ─────────────────────────────────────────────────────────────────────────────

// Create handles POST /api/v1/items
func (h *ItemHandler) Create(c *gin.Context) {
	// 1+2. Parse & validate
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	// 3. Call service (DTO → service input via ToService())
	item, err := h.itemSvc.CreateItem(c.Request.Context(), req.ToService())
	if err != nil {
		respondError(c, err)
		return
	}

	// 4. Respond (domain → response DTO)
	respondCreated(c, fmt.Sprintf("آیتم %s با موفقیت اضافه شد", item.Name))
}

// GetByID handles GET /api/v1/items/:id
func (h *ItemHandler) GetByID(c *gin.Context) {
	// 1. Parse path param
	id, err := parseID(c, "id")
	if err != nil {
		return // parseID already wrote the error response
	}

	// 3. Call service
	item, err := h.itemSvc.GetItem(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	// 4. Respond
	respondOK(c, dto.ItemFromDomain(item))
}

// Update handles PATCH /api/v1/items/:id
// Uses PATCH (not PUT) so clients can send only the fields they want to change.
func (h *ItemHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	item, err := h.itemSvc.UpdateItem(c.Request.Context(), id, req.ToService())
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ItemFromDomain(item))
}

// Delete handles DELETE /api/v1/items/:id
func (h *ItemHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	if err := h.itemSvc.DeleteItem(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	respondNoContent(c)
}

// List handles GET /api/v1/items
// Supports optional query params: ?category_id=1&is_available=true&page=1&page_size=20
func (h *ItemHandler) List(c *gin.Context) {
	// 1. Parse query params into filter
	filter := port.ItemFilter{
		CategoryId: c.GetInt64("categoryId"),
	}

	// 3. Call service
	items, total, err := h.itemSvc.ListItems(c.Request.Context(), filter)
	if err != nil {
		respondError(c, err)
		return
	}

	// 4. Convert each domain item → DTO
	result := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		result[i] = dto.ItemFromDomain(item)
	}

	respondList(c, result, total)
}

// SetAvailability handles PATCH /api/v1/items/:id/availability
// A dedicated endpoint because toggling availability is a frequent operation.
func (h *ItemHandler) SetAvailability(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		return
	}

	var req dto.SetAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	item, err := h.itemSvc.SetAvailability(c.Request.Context(), id, req.IsAvailable)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ItemFromDomain(item))
}

// ── Private helpers ───────────────────────────────────────────────────────────

// parseID extracts an int64 path param and writes the error response if it fails.
// Handlers call:  id, err := parseID(c, "id"); if err != nil { return }
func parseID(c *gin.Context, param string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorBody{
			Error:   "bad_request",
			Message: "invalid " + param + " — must be an integer",
		})
	}
	return id, err
}

// queryInt reads a query param as int with a fallback default.
func queryInt(c *gin.Context, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultVal
	}
	return n
}
