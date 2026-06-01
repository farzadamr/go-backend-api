package dependency

import (
	"github.com/farzadamr/go-backend-api/internal/infra/database"
	"github.com/farzadamr/go-backend-api/internal/infra/repository"
	"github.com/farzadamr/go-backend-api/internal/port"
)

func GetItemRepository() port.ItemRepository {
	var preloads []database.PreloadEntity = []database.PreloadEntity{{Entity: "File"}, {Entity: "Category"}}
	return repository.NewItemRepository(preloads)
}
