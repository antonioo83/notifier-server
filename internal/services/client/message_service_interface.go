package client

type MessageService interface {
	SendMessages(filepath string) (status int, err error)
	GetStatus(messageId string) (MessageStatusResponse, error)
}
