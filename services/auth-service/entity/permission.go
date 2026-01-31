package entity

import (
	"errors"
)

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PaginationMeta struct {
	Page        int64 `json:"page"`
	PerPage     int64 `json:"perPage"`
	Total       int64 `json:"total"`
	TotalItems  int64 `json:"total_items"`
	HasNextPage bool  `json:"has_next_page"`
	HasPrevPage bool  `json:"has_prev_page"`
}

func NewPermission(name string) (*Permission, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	return &Permission{
		ID:   1,
		Name: name,
	}, nil
}
func (p *Permission) id() int {
	return p.ID
}
func (p *Permission) name() string {
	return p.Name
}
