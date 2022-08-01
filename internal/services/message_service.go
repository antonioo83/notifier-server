package services

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"hash/fnv"
	"time"
)

type MessageCreateRequest struct {
	MessageId               string `validate:"required,max=64" faker:"uuid_hyphenated" json:"messageId,omitempty`
	Priority                string `validate:"required,oneof='high' 'low' 'normal'" faker:"oneof: high,low,normal" json:"priority,omitempty"`
	URL                     string `validate:"required,max=1000" faker:"url" json:"url,omitempty"`
	Command                 string `validate:"required,max=10" faker:"oneof: post,put" json:"command,omitempty"`
	Content                 string `validate:"required" faker:"len=20" json:"content,omitempty"`
	SendAt                  string `validate:"-" faker:"timestamp" json:"sendAt,omitempty"`
	SuccessHttpStatus       int    `validate:"numeric" faker:"oneof: 200,201,204" json:"successHttpStatus,omitempty"`
	SuccessResponse         string `validate:"max=300" faker:"len=10" json:"successResponse,omitempty"`
	Description             string `validate:"max=100" faker:"len=100" json:"description,omitempty"`
	isSendNotReceivedNotify bool
}

const (
	ErrorCode   = 1
	SuccessCode = 0
)

type MessageCreateResponse struct {
	MessageID string
	Code      int
	Message   string
	err       error
}

type MessageRouteParameters struct {
	Config             config.Config
	UserRepository     interfaces.UserRepository
	ResourceRepository interfaces.ResourceRepository
	MessageRepository  interfaces.MessageRepository
}

func CreateMessages(userAuth auth.UserAuth, messageRequests []MessageCreateRequest, param MessageRouteParameters) ([]MessageCreateResponse, error) {
	var responses []MessageCreateResponse
	validate := validator.New()
	for _, messageRequest := range messageRequests {
		response := MessageCreateResponse{MessageID: messageRequest.MessageId, Code: ErrorCode}
		err := validate.Struct(messageRequest)
		if err != nil {
			return nil, fmt.Errorf("this request has mistake: %w", err)
		}

		isExist, err := param.MessageRepository.IsInDatabase(messageRequest.MessageId)
		if err != nil {
			return nil, fmt.Errorf("can't get user from the database: %w", err)
		}
		if isExist {
			return nil, fmt.Errorf("this message already exist: %w", err)
		}

		hash := getHash(messageRequest.URL)
		resource, err := param.ResourceRepository.FindByCode(int(hash))
		if err != nil {
			return nil, fmt.Errorf("can't get resource from the database: %w", err)
		}
		resourceId := 0
		if resource == nil {
			resourceId, err = param.ResourceRepository.Save(models.Resource{
				UserId: userAuth.User.ID,
				Code:   int64(hash),
				URL:    messageRequest.URL,
			})
			if err != nil {
				return nil, fmt.Errorf("can't create resource in the database: %w", err)
			}
		} else {
			resourceId = resource.ID
		}

		var message models.Message
		err = copier.Copy(&message, &messageRequest)
		message.UserId = userAuth.User.ID
		message.ResourceId = resourceId
		if err != nil {
			return nil, fmt.Errorf("can't copy data from the request: %w", err)
		}
		if messageRequest.SendAt != "" {
			message.SendAt, err = getTimeFromStr(messageRequest.SendAt)
			if err != nil {
				return nil, fmt.Errorf("request has wrong sentAt value: %w", err)
			}
		}

		err = param.MessageRepository.Save(message)
		if err != nil {
			return nil, fmt.Errorf("can't create a message in the database: %w", err)
		}

		response.Code = SuccessCode
		response.Message = "Ok"
		responses = append(responses, response)
	}

	return responses, nil
}

func getHash(s string) uint32 {
	h := fnv.New32()
	h.Write([]byte(s))
	return h.Sum32()
}

type MessageDeleteRequest struct {
	MessageId string `validate:"required,max=64"`
}

func DeleteMessage(r MessageDeleteRequest, param MessageRouteParameters) error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return fmt.Errorf("this request has mistake: %w", err)
	}

	isExist, err := param.MessageRepository.IsInDatabase(r.MessageId)
	if err != nil {
		return fmt.Errorf("can't get message from the database: %w", err)
	}
	if !isExist {
		return fmt.Errorf("this message isn't exist, userId=%s", r.MessageId)
	}

	err = param.MessageRepository.Delete(r.MessageId)
	if err != nil {
		return fmt.Errorf("can't delete an message from the database: %w", err)
	}

	return nil
}

type MessageGetRequest struct {
	MessageId string `validate:"required,min=1,max=64"`
}

type MessageResponse struct {
	MessageId               string `json:"messageId`
	Priority                string `json:"priority"`
	URL                     string `json:"url"`
	Command                 string `json:"command"`
	Content                 string `json:"content"`
	SendAt                  string `json:"sendAt"`
	SuccessHttpStatus       int    `json:"successHttpStatus"`
	SuccessResponse         string `json:"successResponse"`
	Description             string `json:"description"`
	IsSendNotReceivedNotify bool   `json:"isSendNotReceivedNotify"`
	IsSent                  bool   `json:"isSent"`
	AttemptCount            int    `json:"attemptCount"`
	IsSentCallback          bool   `json:"isSentCallback"`
	CallbackAttemptCount    int    `json:"callbackAttemptCount"`
	CreatedAt               string `json:"createdAt"`
}

func GetMessage(httpRequest MessageGetRequest, param MessageRouteParameters) (*MessageResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	message, err := param.MessageRepository.FindByCode(httpRequest.MessageId)
	if err != nil {
		return nil, fmt.Errorf("can't get an message from the database: %w", err)
	}

	if message == nil {
		return nil, NotFoundError
	}

	var response MessageResponse
	err = copier.Copy(&response, &message)
	if err != nil {
		return nil, fmt.Errorf("can't copy data for the response: %w", err)
	}
	response.SendAt = message.SendAt.Format("2006-01-02 15:04:05")
	response.CreatedAt = message.CreatedAt.Format("2006-01-02 15:04:05")

	response.URL = message.Resource.URL

	return &response, nil
}

func getTimeFromStr(dateTime string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	time, err := time.Parse(layout, dateTime)
	if err != nil {
		return time, err
	}

	return time, nil
}
