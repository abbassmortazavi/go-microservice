package entity

import (
	"errors"
)

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func NewPermission(name, description string) (*Permission, error) {
	if name == "" || description == "" {
		return nil, errors.New("name or description is empty")
	}
	return &Permission{
		ID:   1,
		Name: name,
		Desc: description,
	}, nil
}
func (p *Permission) id() int {
	return p.ID
}
func (p *Permission) name() string {
	return p.Name
}
func (p *Permission) description() string {
	return p.Desc
}
