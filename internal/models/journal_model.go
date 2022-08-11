package models

import "time"

type Journal struct {
	ID              int       `copier:"-"`               // Unique identification of a record.
	UserId          int       `copier:"-"`               // User ID.
	ResourceId      int       `copier:"-"`               // External resource ID.
	MessageId       int       `copier:"-"`               // Message ID.
	ResponseStatus  int       `copier:"ResponseStatus"`  // Received HTTP status from a resource.
	ResponseContent string    `copier:"ResponseContent"` // Received content from a resource.
	Description     string    `copier:"Description"`     // Additional comment for a response of a resource.
	CreatedAt       time.Time // Data and Time of the created record.
}
