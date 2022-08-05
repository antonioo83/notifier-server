package grpc

import (
	"context"
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/server/grpc/user_proto"
	"github.com/antonioo83/notifier-server/internal/services"
	"github.com/jinzhu/copier"
	"net/http"
)

// UserServer поддерживает все необходимые методы сервера.
type UserServer struct {
	user_proto.UnimplementedUserServer
	Config         config.Config
	UserRepository interfaces.UserRepository
}

func (s *UserServer) Create(ctx context.Context, in *user_proto.UserCreateRequest) (*user_proto.UserCreateResponse, error) {
	response := user_proto.UserCreateResponse{}

	var req services.UserCreateRequest
	err := copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}

	var param services.UserRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	token, err := services.CreateUser(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't create an user: %w", err)
	}

	response = user_proto.UserCreateResponse{
		Status:  http.StatusCreated,
		Message: "user was saved",
		Token:   token,
	}

	return &response, err
}

func (s *UserServer) Update(ctx context.Context, in *user_proto.UserUpdateRequest) (*user_proto.UserUpdateResponse, error) {
	response := user_proto.UserUpdateResponse{}
	user, err := s.UserRepository.FindByCode(in.UserId)
	if err != nil || user == nil {
		return &response, fmt.Errorf("i can't find user: %w", err)
	}

	var req services.UserCreateRequest
	err = copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}

	var param services.UserRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	token, err := services.UpdateUser(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't create an user: %w", err)
	}

	response = user_proto.UserUpdateResponse{
		Status:  http.StatusNoContent,
		Message: "user was updated",
		Token:   token,
	}

	return &response, err
}

func (s *UserServer) Delete(ctx context.Context, in *user_proto.UserDeleteRequest) (*user_proto.UserDeleteResponse, error) {
	response := user_proto.UserDeleteResponse{}
	user, err := s.UserRepository.FindByCode(in.UserId)
	if err != nil || user == nil {
		return &response, fmt.Errorf("i can't find user: %w", err)
	}

	var req services.UserDeleteRequest
	err = copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}

	var param services.UserRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository
	err = services.DeleteUser(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't create an user: %w", err)
	}

	response = user_proto.UserDeleteResponse{
		Status:  http.StatusAccepted,
		Message: "user was deleted",
	}

	return &response, err
}

func (s *UserServer) GetUser(ctx context.Context, in *user_proto.UserGetRequest) (*user_proto.UserGetResponse, error) {
	response := user_proto.UserGetResponse{}
	var req services.UserGetRequest
	err := copier.Copy(&req, &in)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}
	var param services.UserRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository

	result, err := services.GetUser(req, param)
	if err != nil {
		return &response, fmt.Errorf("i can't get short url: %w", err)
	}

	err = copier.Copy(&response, &result)
	if err != nil {
		return &response, fmt.Errorf("i can't copy data: %w", err)
	}
	response.UserId = result.UserId

	return &response, nil
}

func (s *UserServer) GetUsers(ctx context.Context, in *user_proto.UsersGetRequest) (*user_proto.UserGetResponses, error) {
	responses := user_proto.UserGetResponses{}
	var req services.UsersGetRequest
	err := copier.Copy(&req, &in)
	if err != nil {
		return &responses, fmt.Errorf("i can't copy data: %w", err)
	}
	var param services.UserRouteParameters
	param.Config = s.Config
	param.UserRepository = s.UserRepository

	results, err := services.GetUsers(req, param)
	if err != nil {
		return &responses, fmt.Errorf("i can't get short url: %w", err)
	}

	for _, result := range results {
		response := user_proto.UserGetResponse{}
		err = copier.Copy(&response, &result)
		if err != nil {
			return &responses, fmt.Errorf("i can't copy data: %w", err)
		}
		response.UserId = result.UserId
		responses.Items = append(responses.Items, &response)
	}

	return &responses, nil
}
