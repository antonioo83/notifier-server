package services

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"strings"
)

type UserRequest struct {
	UserId            string `validate:"required,max=64"`
	Role              string `validate:"required,oneof='service' 'device'"`
	Title             string `validate:"required,max=100"`
	Description       string `validate:"max=256"`
	IsRegenerateToken bool
}

type UserRouteParameters struct {
	Config         config.Config
	UserRepository interfaces.UserRepository
}

var RequestError = fmt.Errorf("request has wrong data")

func CreateUser(userRequest UserRequest, param UserRouteParameters) (token string, error error) {
	validate := validator.New()
	err := validate.Struct(userRequest)
	if err != nil {
		return "", RequestError
	}

	isExist, err := param.UserRepository.IsInDatabase(userRequest.UserId)
	if err != nil {
		return "", fmt.Errorf("can't get user from the database: %w", err)
	}
	if isExist {
		return "", fmt.Errorf("this user already is exist, userId=%s", userRequest.UserId)
	}

	authToken, err := getAuthToken()
	if err != nil {
		return "", fmt.Errorf("can't generate user auth token: %w", err)
	}

	var user models.User
	err = copier.Copy(&user, &userRequest)
	if err != nil {
		return "", fmt.Errorf("can't copy data from the request: %w", err)
	}
	user.AuthToken = authToken
	err = param.UserRepository.Save(user)
	if err != nil {
		return "", fmt.Errorf("can't create an user in the database: %w", err)
	}

	return authToken, nil
}

func UpdateUser(userRequest UserRequest, param UserRouteParameters) (token string, error error) {
	validate := validator.New()
	err := validate.Struct(userRequest)
	if err != nil {
		return "", RequestError
	}

	model, err := param.UserRepository.FindByCode(userRequest.UserId)
	if err != nil {
		return "", fmt.Errorf("can't get user from the database: %w", err)
	}
	if model == nil {
		return "", fmt.Errorf("this user isn't exist, userId=%s", userRequest.UserId)
	}

	authToken := model.AuthToken
	if userRequest.IsRegenerateToken {
		authToken, err = getAuthToken()
		if err != nil {
			return "", fmt.Errorf("can't generate user auth token: %w", err)
		}
	}

	var user models.User
	err = copier.Copy(&user, &userRequest)
	if err != nil {
		return "", fmt.Errorf("can't copy data from the request: %w", err)
	}
	user.ID = model.ID
	user.AuthToken = authToken
	err = param.UserRepository.Update(user)
	if err != nil {
		return "", fmt.Errorf("can't update an user in the database: %w", err)
	}

	return authToken, nil
}

func getAuthToken() (string, error) {
	uuidWithHyphen, err := uuid.NewRandom()
	if err != nil {
		return "", nil
	}
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	return uuid, nil
}

type UserDeleteRequest struct {
	UserId string `validate:"required,max=64"`
}

func DeleteUser(userDeleteRequest UserDeleteRequest, param UserRouteParameters) error {
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

type UserGetRequest struct {
	UserId string `validate:"required,min=1,max=64"`
}

type UserResponse struct {
	UserId      string `json:"userId,omitempty" copier:"Code"`
	Role        string `json:"role,omitempty" copier:"Role"`
	Title       string `json:"title,omitempty" copier:"Title"`
	Description string `json:"description,omitempty" copier:"Description"`
}

func GetUser(httpRequest UserGetRequest, param UserRouteParameters) (*UserResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	user, err := param.UserRepository.FindByCode(httpRequest.UserId)
	if err != nil {
		return nil, fmt.Errorf("can't find an user in the database: %w", err)
	}

	var response UserResponse
	err = copier.Copy(&response, &user)
	if err != nil {
		return nil, fmt.Errorf("can't copy data for the response: %w", err)
	}

	return &response, nil
}

type UsersGetRequest struct {
	Limit  int `validate:"numeric"`
	Offset int `validate:"numeric"`
}

func GetUsers(httpRequest UsersGetRequest, param UserRouteParameters) ([]UserResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	users, err := param.UserRepository.FindAll(1000, 0)
	if err != nil {
		return nil, fmt.Errorf("can't find users in the database: %w", err)
	}

	responses, err := getUsersResponses(users)
	if err != nil {
		return nil, fmt.Errorf("can't get responses: %w", err)
	}

	return responses, nil
}

func getUsersResponses(users *map[int]models.User) ([]UserResponse, error) {
	var responses []UserResponse
	for _, user := range *users {
		var response UserResponse
		err := copier.Copy(&response, &user)
		if err != nil {
			return responses, fmt.Errorf("can't copy data for the response: %w", err)
		}
		responses = append(responses, response)
	}

	return responses, nil
}
