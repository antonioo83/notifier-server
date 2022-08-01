package models

import "time"

type Message struct {
	ID                   int    `copier:"-"`
	Code                 string `copier:"MessageId"`
	UserId               int    `copier:"-"`
	ResourceId           int    `copier:"-"`
	Command              string `copier:"Command"`
	Priority             string `copier:"Priority"`
	Content              string `copier:"Content"`
	IsSent               bool   `copier:"IsSent"`
	AttemptCount         int    `copier:"AttemptCount"`
	IsSentCallback       bool   `copier:"IsSentCallback"`
	CallbackAttemptCount int    `copier:"CallbackAttemptCount"`
	SuccessHttpStatus    int    `copier:"SuccessHttpStatus"`
	SuccessResponse      string `copier:"SuccessResponse"`
	Description          string `copier:"Description"`
	SendAt               time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	Resource             Resource
}
