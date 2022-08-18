package services

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

type SettingCreateRequest struct {
	ResourceId  string `validate:"required,max=64" faker:"uuid_hyphenated" json:"settingId,omitempty`
	Title       string `validate:"max=100" faker:"username" json:"title,omitempty"`
	CallbackURL string `validate:"max=1000" faker:"oneof: service,device" json:"role,omitempty"`
	Count       int    `validate:"-" faker:"oneof: service,device" json:"count,omitempty"`
	Intervals   []int  `validate:"-" faker:"oneof: service,device" json:"intervals,omitempty"`
	Timeout     int    `validate:"-" faker:"oneof: service,device" json:"timeout,omitempty"`
	Description string `validate:"max=256" faker:"len=256" json:"description,omitempty"`
}

type SettingRouteParameters struct {
	Config            config.Config
	SettingRepository interfaces.SettingRepository
}

// CreateSetting create an user by request.
func CreateSetting(request SettingCreateRequest, param SettingRouteParameters) (error error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return RequestError
	}

	isExist, err := param.SettingRepository.IsInDatabase(request.ResourceId)
	if err != nil {
		return fmt.Errorf("can't get settings from the database: %w", err)
	}
	if isExist {
		return fmt.Errorf("this settings already is exist, settingId=%s", request.ResourceId)
	}

	var setting models.Setting
	err = copier.Copy(&setting, &request)
	if err != nil {
		return fmt.Errorf("can't copy data from the request: %w", err)
	}
	err = param.SettingRepository.Save(setting)
	if err != nil {
		return fmt.Errorf("can't create an setting in the database: %w", err)
	}

	return nil
}

// UpdateSetting update an user by request.
func UpdateSetting(request SettingCreateRequest, param SettingRouteParameters) (error error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return RequestError
	}

	model, err := param.SettingRepository.FindByCode(request.ResourceId)
	if err != nil {
		return fmt.Errorf("can't get setting from the database: %w", err)
	}
	if model == nil {
		return fmt.Errorf("this setting isn't exist, userId=%s", request.ResourceId)
	}

	var setting models.Setting
	err = copier.Copy(&setting, &request)
	if err != nil {
		return fmt.Errorf("can't copy data from the request: %w", err)
	}
	err = param.SettingRepository.Update(setting)
	if err != nil {
		return fmt.Errorf("can't update an setting in the database: %w", err)
	}

	return nil
}

type SettingDeleteRequest struct {
	ResourceId string `validate:"required,max=64"`
}

// DeleteSetting delete an user by request.
func DeleteSetting(settingDeleteRequest SettingDeleteRequest, param SettingRouteParameters) error {
	validate := validator.New()
	err := validate.Struct(settingDeleteRequest)
	if err != nil {
		return RequestError
	}

	isExist, err := param.SettingRepository.IsInDatabase(settingDeleteRequest.ResourceId)
	if err != nil {
		return fmt.Errorf("can't get user from the database: %w", err)
	}
	if !isExist {
		return fmt.Errorf("this user isn't exist, userId=%s", settingDeleteRequest.ResourceId)
	}

	err = param.SettingRepository.Delete(settingDeleteRequest.ResourceId)
	if err != nil {
		return fmt.Errorf("can't delete an user from the database: %w", err)
	}

	return nil
}

type SettingGetRequest struct {
	UserId string `validate:"required,min=1,max=64"`
}

type SettingResponse struct {
	UserId      string `json:"userId,omitempty" copier:"Code"`
	Role        string `json:"role,omitempty" copier:"Role"`
	Title       string `json:"title,omitempty" copier:"Title"`
	Description string `json:"description,omitempty" copier:"Description"`
}

// GetUser get an user by request.
func GetSetting(httpRequest SettingGetRequest, param UserRouteParameters) (*UserResponse, error) {
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

type UsersGetRequest struct {
	Limit  int `validate:"numeric"`
	Offset int `validate:"numeric"`
}

// GetUsers get users by request.
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

// getUsersResponses gets user responses.
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
