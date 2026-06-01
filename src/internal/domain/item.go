package domain

import "time"

type Item struct {
	Id          int64
	Name        string
	FileId      int64
	Description string
	Price       int
	IsAvailable bool
	Ingredients string
	CategoryId  int64
	CreateAt    time.Time
	UpdateAt    time.Time
}

func (i *Item) IsOrderable() bool {
	return i.IsAvailable && i.Price > 0
}
