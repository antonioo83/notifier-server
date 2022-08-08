package models

import "time"

type User struct {
	ID          int       `copier:"-"`           // User ID.
	Code        string    `copier:"UserId"`      // Unique identification of a user in the external source.
	Role        string    `copier:"Role"`        // Role of the user admin or external source.
	Title       string    `copier:"Title"`       // Name of the external source.
	AuthToken   string    `copier:"-"`           // Token for the authentication.
	Description string    `copier:"Description"` // Additional comment for an user.
	CreatedAt   time.Time // Date and time the message was created.
	UpdatedAt   time.Time // Date and time the message was updated.
	DeletedAt   time.Time // Date and time the message was deleted.
}
