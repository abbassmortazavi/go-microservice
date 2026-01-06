package entity

import (
	"errors"
)

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewPermission(name string) (*Permission, error) {
	if name == "" {
		return nil, errors.New("name or description is empty")
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
