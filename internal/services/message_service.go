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
)

type MessageCreateRequest struct {
	MessageId               string `validate:"required,max=64" faker:"uuid_hyphenated" json:"messageId,omitempty`
	Priority                string `validate:"required,oneof='high' 'low' 'normal'" faker:"oneof: high,low,normal" json:"priority,omitempty"`
	URL                     string `validate:"required,max=1000" faker:"url" json:"url,omitempty"`
	Command                 string `validate:"required,max=10" faker:"oneof: post,put" json:"command,omitempty"`
	Content                 string `validate:"required" faker:"len=20" json:"title,omitempty"`
	SendAt                  string `validate:"datetime" faker:"timestamp" json:"sendAt,omitempty"`
	SuccessHttpStatus       int    `validate:"numeric" faker:"oneof: 200,201,204" json:"successHttpStatus,omitempty"`
	SuccessResponse         string `validate:"max=300" faker:"len=10" json:"successResponse,omitempty"`
	Description             string `validate:"max=100" faker:"len=256" json:"description,omitempty"`
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
			response.Message = err.Error()
			responses = append(responses, response)
			continue
		}

		isExist, err := param.MessageRepository.IsInDatabase(messageRequest.MessageId)
		if err != nil {
			return nil, fmt.Errorf("can't get user from the database: %w", err)
		}
		if isExist {
			response.Message = "this user already is exist"
			responses = append(responses, response)
			continue
		}

		hash := getHash(messageRequest.URL)
		resource, err := param.ResourceRepository.FindByCode(int(hash))
		if err != nil {
			return nil, fmt.Errorf("can't get user from the database: %w", err)
		}
		resourceId := 0
		if resource == nil {
			resourceId, err = param.ResourceRepository.Save(models.Resource{
				UserId: userAuth.User.ID,
				Code:   int(hash),
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

/*
type MessageDeleteRequest struct {
	UserId string `validate:"required,max=64"`
}

func DeleteMessage(userDeleteRequest MessageDeleteRequest, param MessageRouteParameters) error {
	validate := validator.New()
	err := validate.Struct(userDeleteRequest)
	if err != nil {
		return RequestError
	}

	isExist, err := param.UserRepository.IsInDatabase(userDeleteRequest.UserId)
	if err != nil {
		return fmt.Errorf("can't get user from the database: %w", err)
	}
	if !isExist {
		return fmt.Errorf("this user isn't exist, userId=%s", userDeleteRequest.UserId)
	}

	err = param.UserRepository.Delete(userDeleteRequest.UserId)
	if err != nil {
		return fmt.Errorf("can't delete an user from the database: %w", err)
	}

	return nil
}

type MessageGetRequest struct {
	UserId string `validate:"required,min=1,max=64"`
}

type MessageResponse struct {
	UserId      string `json:"userId,omitempty" copier:"Code"`
	Role        string `json:"role,omitempty" copier:"Role"`
	Title       string `json:"title,omitempty" copier:"Title"`
	Description string `json:"description,omitempty" copier:"Description"`
}

func GetMessage(httpRequest MessageGetRequest, param MessageRouteParameters) (*MessageResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	user, err := param.UserRepository.FindByCode(httpRequest.UserId)
	if err != nil {
		return nil, fmt.Errorf("can't get an user from the database: %w", err)
	}

	if user == nil {
		return nil, NotFoundError
	}

	var response UserResponse
	err = copier.Copy(&response, &user)
	if err != nil {
		return nil, fmt.Errorf("can't copy data for the response: %w", err)
	}

	return &response, nil
}
*/
