package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	r "github.com/m-a-r-a-t/go-jwt-auth/internal"
	authModels "github.com/m-a-r-a-t/go-jwt-auth/internal/models"
	my_jwt "github.com/m-a-r-a-t/go-jwt-auth/pkg/jwt"
	"github.com/m-a-r-a-t/go-rest-wrap/pkg/middleware"
	"github.com/m-a-r-a-t/go-rest-wrap/pkg/router"
)

type AuthControllerServices struct {
	AuthService IAuthService
}

func init() {
	u := r.R.SubRouter("/auth")
	fmt.Println(u.BasePath)
	authRepos := AuthControllerServices{
		AuthService: &r.AuthService,
	}

	tokensConfig := authModels.TokensConfig{
		AccessTokenDuration:  10 * time.Minute,  // 10 minutes
		RefreshTokenDuration: 1440 * time.Hour,  // 60 days
		SignForAccessToken:   []byte("access"),  // подпись для access token
		SignForRefreshToken:  []byte("refresh"), // подпись для refresh token
	}
	// u.GET("/get_user", router.Params{StructureForValidateQueryDataCreator: NewPerson},
	// 	func(ctx router.HandlerCtx) (router.Result, router.Error) {
	// 		p := (*ctx.Structure).(*Person)
	// 		fmt.Println("Name", *p.Name)
	// 		return router.Result{Body: p}, nil
	// 	})

	u.GET("/get_user", router.Params{ // ! тестовый роут с проверкой токена
		Middlewares: middleware.New(r.VerifyAccessToken),
	},
		func(ctx router.HandlerCtx) (router.Result, router.Error) {
			return router.Result{Body: "Available route"}, nil
		})

	// TODO ! сделать /refresh роут , рефрешим оба токена

	/*

		Общий логин по сервисам YANDEX,VK и тд
	*/
	u.POST("/login", router.Params{
		StructureForValidateQueryDataCreator: authModels.NewLoginRouteValidation,
	},
		func(ctx router.HandlerCtx) (router.Result, router.Error) {
			fmt.Println(ctx.Request.RemoteAddr)

			loginData := (*ctx.Structure).(*authModels.LoginRouteValidation)

			userDataFromService, err := authRepos.AuthService.GetUserDataFromForeignService(
				loginData.OauthToken,
				loginData.ServiceName,
			)

			if err != nil {
				return router.Result{}, NewForeignServiceError(err.Error())
			}

			availableUserData, _ := authRepos.AuthService.GetAvailableUserData(userDataFromService.Email)

			if availableUserData == nil {
				fmt.Println("Нет данных о пользователе")
				err = authRepos.AuthService.CreateUser(userDataFromService)
				
				if err != nil {
					fmt.Println(err)
					return router.Result{}, NewCreateNewUserError(err.Error())
				}

				availableUserData, err = authRepos.AuthService.GetAvailableUserData(userDataFromService.Email)
				if err != nil {
					fmt.Println(err)
					return router.Result{}, err.Error()
				}
			}

			accessToken, err := authRepos.AuthService.CreateSignedToken(
				tokensConfig.SignForAccessToken,
				tokensConfig.AccessTokenDuration,
				availableUserData,
			)

			if err != nil {
				fmt.Println("access token creating error")
				return router.Result{}, err.Error()
			}

			refreshToken, err := authRepos.AuthService.CreateSignedToken(
				tokensConfig.SignForRefreshToken,
				tokensConfig.RefreshTokenDuration,
				availableUserData,
			)

			if err != nil {
				fmt.Println("refresh token creating error")
				return router.Result{}, err.Error()
			}

			_, err = authRepos.AuthService.CreateOrUpdateRefreshTokenInDb(refreshToken.Token, availableUserData.Id)

			if err != nil {
				return router.Result{}, err.Error()
			}
			fmt.Println("Access token:", accessToken)

			return router.Result{Body: authModels.LoginResponse{
				AccessToken:       accessToken.Token,
				RefreshToken:      refreshToken.Token,
				AccessTokenExpire: accessToken.Expire,
			}}, nil
		})

	u.POST("/refresh_token", router.Params{
		StructureForValidateQueryDataCreator: authModels.NewRefreshTokenRouteValidation,
	},
		func(ctx router.HandlerCtx) (router.Result, router.Error) {
			var accessToken, refreshToken *authModels.MyToken
			refreshTokenData := (*ctx.Structure).(*authModels.RefreshTokenRouteValidation)
			verifyResult, err := my_jwt.VerifyJWT(refreshTokenData.RefreshToken, tokensConfig.SignForRefreshToken)

			if err != nil {
				// ! Ошибка валидации токена либо истек либо липовый
				return router.Result{}, err.Error()
			}

			if claims, ok := verifyResult.Token.Claims.(jwt.MapClaims); ok && verifyResult.Result {
				user_data, err := authRepos.AuthService.GetUserDataFromJWT(claims)
				if err != nil {
					return router.Result{}, err.Error()
				}

				isValid, err := authRepos.AuthService.CompareRefreshTokens(refreshTokenData.RefreshToken, user_data.Id)

				if isValid {
					availableUserData, err := authRepos.AuthService.GetAvailableUserData(user_data.Email)

					if err != nil {
						return router.Result{}, err.Error()
					}

					accessToken, err = authRepos.AuthService.CreateSignedToken(
						tokensConfig.SignForAccessToken,
						tokensConfig.AccessTokenDuration,
						availableUserData,
					)

					if err != nil {
						fmt.Println("access token creating error")
						return router.Result{}, err.Error()
					}

					refreshToken, err = authRepos.AuthService.CreateSignedToken(
						tokensConfig.SignForRefreshToken,
						tokensConfig.RefreshTokenDuration,
						availableUserData,
					)

					if err != nil {
						fmt.Println("refresh token creating error")
						return router.Result{}, err.Error()
					}

					_, err = authRepos.AuthService.CreateOrUpdateRefreshTokenInDb(refreshToken.Token, user_data.Id)

					if err != nil {
						return router.Result{}, err.Error()
					}

				} else {
					return router.Result{}, err.Error()
				}

			}

			return router.Result{Body: authModels.LoginResponse{
				AccessToken:       accessToken.Token,
				RefreshToken:      refreshToken.Token,
				AccessTokenExpire: accessToken.Expire,
			}}, nil
		})

}
