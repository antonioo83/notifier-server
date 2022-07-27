package models

import "time"

type User struct {
	ID          int    `copier:"-"`
	Code        string `copier:"UserId"`
	Role        string `copier:"Role"`
	Title       string `copier:"Title"`
	AuthToken   string `copier:"-"`
	Description string `copier:"Description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
