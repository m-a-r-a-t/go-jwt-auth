package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type LoginRouteValidation struct {
	OauthToken  string `json:"oauth_token,omitempty" validate:"required"`
	ServiceName string `json:"service_name,omitempty" validate:"required"`
}

func NewLoginRouteValidation() interface{} {
	return &LoginRouteValidation{}
}

type RefreshTokenRouteValidation struct {
	RefreshToken string `json:"refresh_token,omitempty" validate:"required"`
}

func NewRefreshTokenRouteValidation() interface{} {
	return &RefreshTokenRouteValidation{}
}

type UserDataFromService struct {
	Email     string   `json:"default_email,omitempty" db:"email"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Sex       string   `json:"sex,omitempty"`
	Emails    []string `json:"emails,omitempty"`
	Birthday  string   `json:"birthday,omitempty"`
}

type TokensConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	SignForAccessToken   []byte
	SignForRefreshToken  []byte
}

type LoginResponse struct {
	AccessToken       string `json:"access_token,omitempty"`
	RefreshToken      string `json:"refresh_token,omitempty"`
	AccessTokenExpire int64  `json:"access_token_expire,omitempty"`
}

type UserData struct {
	Id    []uint8 `json:"id,omitempty" db:"id"`
	Email string  `json:"email,omitempty" validate:"required" `
}

// ! сделать маппер из UserDataFromService в  UserData

func MapFromServiceToOwnData(dataFromService interface{}) *UserDataFromService {
	// ! Маппинг происходит по golang тэгам
	var data = UserDataFromService{}

	bytes, err := json.Marshal(dataFromService)
	if err != nil {
		fmt.Println("Can't serislize", dataFromService)
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Не возможно распарсить json")
	}

	return &data
}

func MapUserDataServiceToUserData(data *UserDataFromService) *UserData {
	return &UserData{
		Email: data.Email,
	}
}

type MyToken struct {
	Token  string
	Expire int64
}
