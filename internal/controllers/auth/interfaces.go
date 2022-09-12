package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	authModels "github.com/m-a-r-a-t/go-jwt-auth/internal/models"
)

type IAuthService interface {
	GetUserDataFromForeignService(oauthToken string, serviceName string) (*authModels.UserDataFromService, error)
	CreateSignedToken(sign []byte, duration time.Duration, data interface{}) (*authModels.MyToken, error)
	GetAvailableUserData(email string) (*authModels.UserData, error)
	CreateUser(data *authModels.UserDataFromService) error
	CompareRefreshTokens(tokenFromRequest string, user_id []uint8) (bool, error)
	GetUserDataFromJWT(claims jwt.MapClaims) (*authModels.UserData, error)
	CreateOrUpdateRefreshTokenInDb(refreshToken string, user_id []uint8) (bool, error)
}
