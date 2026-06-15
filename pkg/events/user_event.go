package events

type UserRegistered struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
