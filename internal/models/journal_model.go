package models

import "time"

type Journal struct {
	ID              int    `copier:"-"`
	UserId          int    `copier:"-"`
	ResourceId      int    `copier:"-"`
	MessageId       int    `copier:"-"`
	ResponseStatus  int    `copier:"ResponseStatus"`
	ResponseContent string `copier:"ResponseContent"`
	Description     string `copier:"Description"`
	CreatedAt       time.Time
}
