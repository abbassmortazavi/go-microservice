package entity

import "errors"

type Role struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

func NewRole(name string) (*Role, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}
	return &Role{
		ID:          2,
		Name:        name,
		Permissions: make([]Permission, 0),
	}, nil
}
