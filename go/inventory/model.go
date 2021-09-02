package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InventoryItem struct {
	ID           string       `json:"id" binding:"required,uuid"`
	Name         string       `json:"name" binding:"required"`
	ReleaseDate  time.Time    `json:"releaseDate" binding:"required"`
	Manufacturer Manufacturer `json:"manufacturer" binding:"required"`
}

type Manufacturer struct {
	Name     string `json:"name" binding:"required"`
	HomePage string `json:"homePage" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type inventory struct {
	items map[string]InventoryItem
	sync.RWMutex
}

func (i *inventory) add(item InventoryItem) error {
	i.Lock()
	defer i.Unlock()
	_, ok := i.items[item.ID]
	if ok {
		return fmt.Errorf("duplicate key %q", item.ID)
	}
	i.items[item.ID] = item
	return nil
}

func (i *inventory) search(s string, skip, limit int) []InventoryItem {
	i.RLock()
	defer i.RUnlock()
	res := make([]InventoryItem, 0)
	for _, i := range i.items {
		if s == "" || strings.Contains(i.Name, s) {
			res = append(res, i)
		}
	}

	if s != "" {
		res = make([]InventoryItem, 0)
		for _, i := range i.items {
			if strings.Contains(i.Name, s) {
				res = append(res, i)
			}
		}
	}
	if skip > 0 && skip < len(res) {
		res = res[skip:]
	}
	if limit > 0 && limit < len(res) {
		res = res[:limit]
	}
	return res
}

func newInventory() (*inventory, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	data := []InventoryItem{
		{
			ID:          uuid.NewString(),
			Name:        "Surface Laptop 4",
			ReleaseDate: time.Date(2021, 4, 13, 5, 0, 0, 0, loc),
			Manufacturer: Manufacturer{
				Name:     "Microsoft",
				HomePage: "https://www.microsoft.com",
				Phone:    "1-800-642-7676 ",
			},
		},
		{
			ID:          uuid.NewString(),
			Name:        "MacBook Air M1",
			ReleaseDate: time.Date(2020, 11, 10, 5, 0, 0, 0, loc),
			Manufacturer: Manufacturer{
				Name:     "Apple",
				HomePage: "https://www.apple.com",
				Phone:    "1-877-233-8552",
			},
		},
	}
	inv := inventory{items: make(map[string]InventoryItem)}
	for _, i := range data {
		inv.items[i.ID] = i
	}
	return &inv, nil
}
