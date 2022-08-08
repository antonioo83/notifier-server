package models

import "time"

type Resource struct {
	ID          int       `copier:"-"`           // External resource ID.
	UserId      int       `copier:"-"`           // User ID.
	Code        int64     `copier:"-"`           // Unique identification of a resource in the external source.
	URL         string    `copier:"URL"`         // URL of the resource where the request should be sent.
	Description string    `copier:"Description"` // Additional comment for a resource.
	CreatedAt   time.Time // Date and time the message was created.
	UpdatedAt   time.Time // Date and time the message was updated.
	DeletedAt   time.Time // Date and time the message was deleted.
}
