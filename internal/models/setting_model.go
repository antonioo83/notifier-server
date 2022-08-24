package models

import "time"

type Setting struct {
	ID          int       `copier:"-"`           // External resource ID.
	Code        string    `copier:"SettingId"`   // Unique identification of a resource in the external source.
	UserId      int       `copier:"-"`           // User ID.
	ResourceId  int       `copier:"-"`           // External resource ID.
	Title       string    `copier:"Title"`       // Additional comment for a resource.
	Count       int       `copier:"Count"`       // Count of attempt for send to the resource.
	Intervals   []int     `copier:"Intervals"`   // Count of attempt for send to the resource.
	Timeout     int       `copier:"Timeout"`     // URL of the resource where the request should be sent.
	CallbackURL string    `copier:"CallbackURL"` // URL of the resource where the request should be sent.
	Description string    `copier:"Description"` // Additional comment for a resource.
	CreatedAt   time.Time // Date and time the message was created.
	UpdatedAt   time.Time // Date and time the message was updated.
	DeletedAt   time.Time // Date and time the message was deleted.
	User        User      // User relation.
	Resource    Resource  // Resource relation.
}
