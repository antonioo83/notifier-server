package grpc

import (
	"context"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/server/grpc/message_proto"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/jinzhu/copier"
	"net/http"
)

type MessageServer struct {
	message_proto.UnimplementedMessageServer
	Config             config.Config
	UserAuth           auth.UserAuthService
	UserRepository     interfaces.UserRepository
	ResourceRepository interfaces.ResourceRepository
	MessageRepository  interfaces.MessageRepository
	JournalRepository  interfaces.JournalRepository
}

// CreateMessages create messages by grpc service.
func (s *MessageServer) CreateMessages(ctx context.Context, in *message_proto.MessagesCreateRequest) (*message_proto.MessagesCreateResponse, error) {
	response := message_proto.MessagesCreateResponse{}

	userAuth, err := s.UserAuth.GetAuthUser(in.UserToken)
	if err != nil {
		return &response, fmt.Errorf("can't get an user: %w", err)
	}

	var requests []services.MessageCreateRequest
	for _, item := range in.Items {
		var req services.MessageCreateRequest
		err := copier.Copy(&req, &item)
		if err != nil {
			return &response, fmt.Errorf("i can't copy data: %w", err)
		}
		requests = append(requests, req)
	}

	var param services.MessageRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	param.ResourceRepository = s.ResourceRepository
	param.MessageRepository = s.MessageRepository
	_, err = services.CreateMessages(*userAuth, requests, param)
	if err != nil {
		return &response, fmt.Errorf("i can't create messages: %w", err)
	}

	response = message_proto.MessagesCreateResponse{
		Status:  http.StatusCreated,
		Message: "messages was saved",
	}

	return &response, err
}

// DeleteMessage delete a message by grpc service.
func (s *MessageServer) DeleteMessage(ctx context.Context, in *message_proto.MessageDeleteRequest) (*message_proto.MessageDeleteResponse, error) {
	response := message_proto.MessageDeleteResponse{}
	_, err := s.UserAuth.GetAuthUser(in.UserToken)
	if err != nil {
		return &response, fmt.Errorf("can't get an user: %w", err)
	}

	var req services.MessageDeleteRequest
	err = copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}

	var param services.MessageRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	param.ResourceRepository = s.ResourceRepository
	param.MessageRepository = s.MessageRepository
	err = services.DeleteMessage(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't create an user: %w", err)
	}

	response = message_proto.MessageDeleteResponse{
		Status:  http.StatusAccepted,
		Message: "message was deleted",
	}

	return &response, err
}

// GetMessage get a message by grpc service.
func (s *MessageServer) GetMessage(ctx context.Context, in *message_proto.MessageGetRequest) (*message_proto.MessageGetResponse, error) {
	response := message_proto.MessageGetResponse{}
	var req services.MessageGetRequest
	err := copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}
	var param services.MessageRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	param.ResourceRepository = s.ResourceRepository
	param.MessageRepository = s.MessageRepository
	result, err := services.GetMessage(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't get message: %w", err)
	}

	err = copier.Copy(&response, &result)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}

	return &response, nil
}
