package services

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	"github.com/antonioo83/notifier-server/internal/services/auth"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

type SettingCreateRequest struct {
	SettingId   string `validate:"required,max=64" faker:"uuid_hyphenated" json:"settingId,omitempty"`
	Title       string `validate:"max=100" faker:"len=100" json:"title,omitempty"`
	CallbackURL string `validate:"max=1000" faker:"url" json:"callbackUrl,omitempty"`
	Count       int    `validate:"-" faker:"boundary_start=1, boundary_end=10" json:"count,omitempty"`
	Intervals   []int  `validate:"-" faker:"boundary_start=1, boundary_end=10" json:"intervals,omitempty"`
	Timeout     int    `validate:"-" faker:"boundary_start=1, boundary_end=10" json:"timeout,omitempty"`
	Description string `validate:"max=256" faker:"len=256" json:"description,omitempty"`
}

type SettingRouteParameters struct {
	Config            config.Config
	SettingRepository interfaces.SettingRepository
}

// CreateSetting create an setting by request.
func CreateSetting(userAuth auth.UserAuth, request SettingCreateRequest, param SettingRouteParameters) (error error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return RequestError
	}

	isExist, err := param.SettingRepository.IsInDatabase(userAuth.User.ID, request.SettingId)
	if err != nil {
		return fmt.Errorf("can't get settings from the database: %w", err)
	}
	if isExist {
		return fmt.Errorf("this settings already is exist, settingId=%s", request.SettingId)
	}

	var setting models.Setting
	err = copier.Copy(&setting, &request)
	if err != nil {
		return fmt.Errorf("can't copy data from the request: %w", err)
	}
	setting.UserId = userAuth.User.ID
	err = param.SettingRepository.Save(setting)
	if err != nil {
		return fmt.Errorf("can't create an setting in the database: %w", err)
	}

	return nil
}

// UpdateSetting update an setting by request.
func UpdateSetting(userAuth auth.UserAuth, request SettingCreateRequest, param SettingRouteParameters) (error error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return RequestError
	}

	model, err := param.SettingRepository.FindByCode(userAuth.User.ID, request.SettingId)
	if err != nil {
		return fmt.Errorf("can't get setting from the database: %w", err)
	}
	if model == nil {
		return fmt.Errorf("this setting isn't exist, userId=%s", request.SettingId)
	}

	var setting models.Setting
	err = copier.Copy(&setting, &request)
	if err != nil {
		return fmt.Errorf("can't copy data from the request: %w", err)
	}
	setting.UserId = userAuth.User.ID
	err = param.SettingRepository.Update(setting)
	if err != nil {
		return fmt.Errorf("can't update an setting in the database: %w", err)
	}

	return nil
}

type SettingDeleteRequest struct {
	SettingId string `validate:"required,max=64"`
}

// DeleteSetting delete an setting by request.
func DeleteSetting(userAuth auth.UserAuth, settingDeleteRequest SettingDeleteRequest, param SettingRouteParameters) error {
	validate := validator.New()
	err := validate.Struct(settingDeleteRequest)
	if err != nil {
		return RequestError
	}

	isExist, err := param.SettingRepository.IsInDatabase(userAuth.User.ID, settingDeleteRequest.SettingId)
	if err != nil {
		return fmt.Errorf("can't get setting from the database: %w", err)
	}
	if !isExist {
		return fmt.Errorf("this setting isn't exist, userId=%s", settingDeleteRequest.SettingId)
	}

	err = param.SettingRepository.Delete(userAuth.User.ID, settingDeleteRequest.SettingId)
	if err != nil {
		return fmt.Errorf("can't delete an setting from the database: %w", err)
	}

	return nil
}

type SettingGetRequest struct {
	SettingId string `validate:"required,min=1,max=64"`
}

type SettingResponse struct {
	SettingId   string `json:"settingId,omitempty" copier:"Code"`
	Title       string `json:"title,omitempty" copier:"Title"`
	CallbackURL string `json:"callbackURL,omitempty" copier:"CallbackURL"`
	Count       int    `json:"count,omitempty" copier:"Count"`
	Intervals   []int  `json:"intervals,omitempty" copier:"Intervals"`
	Timeout     int    `json:"timeout,omitempty" copier:"Timeout"`
	Description string `json:"description,omitempty" copier:"Description"`
}

// GetSetting get an setting by request.
func GetSetting(userAuth auth.UserAuth, httpRequest SettingGetRequest, param SettingRouteParameters) (*SettingResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	setting, err := param.SettingRepository.FindByCode(userAuth.User.ID, httpRequest.SettingId)
	if err != nil {
		return nil, fmt.Errorf("can't get an setting from the database: %w", err)
	}

	if setting == nil {
		return nil, NotFoundError
	}

	var response SettingResponse
	err = copier.Copy(&response, &setting)
	if err != nil {
		return nil, fmt.Errorf("can't copy data for the response: %w", err)
	}

	return &response, nil
}

type SettingsGetRequest struct {
	Limit  int `validate:"numeric"`
	Offset int `validate:"numeric"`
}

// GetSettings get users by request.
func GetSettings(userAuth auth.UserAuth, httpRequest SettingsGetRequest, param SettingRouteParameters) ([]SettingResponse, error) {
	validate := validator.New()
	err := validate.Struct(httpRequest)
	if err != nil {
		return nil, RequestError
	}

	users, err := param.SettingRepository.FindAll(userAuth.User.ID, httpRequest.Limit, httpRequest.Offset)
	if err != nil {
		return nil, fmt.Errorf("can't find users in the database: %w", err)
	}

	responses, err := getSettingsResponses(users)
	if err != nil {
		return nil, fmt.Errorf("can't get responses: %w", err)
	}

	return responses, nil
}

// getSettingsResponses gets setting responses.
func getSettingsResponses(users *map[int]models.Setting) ([]SettingResponse, error) {
	var responses []SettingResponse
	for _, setting := range *users {
		var response SettingResponse
		err := copier.Copy(&response, &setting)
		if err != nil {
			return responses, fmt.Errorf("can't copy data for the response: %w", err)
		}
		responses = append(responses, response)
	}

	return responses, nil
}
