package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	authModels "github.com/m-a-r-a-t/go-jwt-auth/internal/models"
	"github.com/m-a-r-a-t/go-jwt-auth/pkg/yandex"
)

type IAuthRepo interface {
	GetUserByEmail(email string) (*authModels.UserData, error)
	CreateUser(userData *authModels.UserData) (bool, error)
	GetUserRefreshTokenById(user_id []uint8) (string, error)
	CreateOrUpdateUserRefreshToken(refreshToken string, user_id []uint8) bool
}

type AuthService struct {
	AuthRepo IAuthRepo
}

func (ar *AuthService) GetUserDataFromForeignService(oauthToken string, serviceName string) (*authModels.UserDataFromService, error) {
	var err error
	var data interface{}
	switch serviceName {

	case "YANDEX":
		data, err = yandex.GetUserFromYandexByOauthToken(oauthToken)
		fmt.Println("YANDEX service")

	case "VK":
		fmt.Println("VK service")

	}
	if err != nil {
		return nil, err
	}

	return authModels.MapFromServiceToOwnData(data), nil

}

/*
Универсальный метод для создания accesss token или refresh token
*/
func (as *AuthService) CreateSignedToken(sign []byte, duration time.Duration, data interface{}) (*authModels.MyToken, error) {
	exp := time.Now().Add(duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  exp,
	})

	tokenString, err := token.SignedString(sign)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// claims["authorized"] = true
	// claims["user"] = "username"
	return &authModels.MyToken{
		Token:  tokenString,
		Expire: exp,
	}, nil
}

func (as *AuthService) GetAvailableUserData(email string) (*authModels.UserData, error) {
	userData, err := as.AuthRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Println("Available user data", err)
		return nil, err
	}
	// запрос к базе
	// обращаемся к другому сервису
	// ! Запрос в базу или на сервис где есть доступ к базе

	return userData, nil
}

func (as *AuthService) CreateUser(data *authModels.UserDataFromService) error {

	result, err := as.AuthRepo.CreateUser(authModels.MapUserDataServiceToUserData(data))
	fmt.Println("Результат создания пользоватля", result)

	if err != nil {
		return err
	}

	return nil

	// запрос к базе
	// обращаемся к другому сервису
	// ! Запрос в базу или на сервис где есть доступ к базе
	// ! использовать env config или любой другой config/config.go например

}

func (as *AuthService) CompareRefreshTokens(tokenFromRequest string, user_id []uint8) (bool, error) {
	tokenFromDb, err := as.AuthRepo.GetUserRefreshTokenById(user_id)

	if err != nil {
		return false, err
	}

	fmt.Println("Token from request", tokenFromRequest)
	fmt.Println("Token from db", tokenFromDb)

	if tokenFromRequest == tokenFromDb {
		return true, nil
	}

	return false, errors.New("tokens not equal")
}

func (as *AuthService) GetUserDataFromJWT(claims jwt.MapClaims) (*authModels.UserData, error) {
	var user_data authModels.UserData
	bytes, err := json.Marshal(claims["data"])
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(bytes, &user_data)

	return &user_data, nil
}

func (as *AuthService) CreateOrUpdateRefreshTokenInDb(refreshToken string, user_id []uint8) (bool, error) {

	result := as.AuthRepo.CreateOrUpdateUserRefreshToken(refreshToken, user_id)

	if result {
		return true, nil
	}

	return false, errors.New("create update token error")
}
