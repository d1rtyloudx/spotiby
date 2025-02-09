package model

import "time"

type Credential struct {
	ID          string
	Username    string
	Email       string
	HashPass    []byte
	Role        string
	IsConfirmed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
