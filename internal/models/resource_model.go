package models

import "time"

type Resource struct {
	ID          int    `copier:"-"`
	UserId      int    `copier:"-"`
	Code        int64  `copier:"-"`
	URL         string `copier:"URL"`
	Description string `copier:"Description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}
