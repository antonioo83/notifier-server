package models

import "time"

type Message struct {
	ID                   int       `copier:"-"`                    // Unique identification of a message.
	Code                 string    `copier:"MessageId"`            // Unique identification of a message in the external source.
	UserId               int       `copier:"-"`                    // User ID.
	ResourceId           int       `copier:"-"`                    // External resource ID.
	Command              string    `copier:"Command"`              // Type of the command(POST or PUT) for send to a resource.
	Priority             string    `copier:"Priority"`             // Priority of a message for send to a resource.
	Content              string    `copier:"Content"`              // Content for send to a resource.
	IsSent               bool      `copier:"IsSent"`               // The marker shows that the message has been sent.
	AttemptCount         int       `copier:"AttemptCount"`         // Count of attempt for send to the resource.
	IsSentCallback       bool      `copier:"IsSentCallback"`       // The marker shows that the callback message has been sent.
	CallbackAttemptCount int       `copier:"CallbackAttemptCount"` // Count of attempt for send to the source.
	SuccessHttpStatus    int       `copier:"SuccessHttpStatus"`    // This status shows that response is success.
	SuccessResponse      string    `copier:"SuccessResponse"`      // This content shows that response is success (not used).
	Description          string    `copier:"Description"`          // Additional comment for a message.
	SendAt               time.Time // At this date and time a message will send.
	CreatedAt            time.Time // Date and time the message was created.
	UpdatedAt            time.Time // Date and time the message was updated.
	DeletedAt            time.Time // Date and time the message was deleted.
	User                 User      // User relation.
	Resource             Resource  // Resource relation.
}
