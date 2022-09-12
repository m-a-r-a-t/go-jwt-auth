package jwt

import (
	"fmt"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v4"
)

/*
Создает мидлварь с определенной подписью в замыкании
*/
func CreateVerifyJWT(s []byte) func(next http.HandlerFunc) http.HandlerFunc {
	sign := s
	return func(next http.HandlerFunc) http.HandlerFunc {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header["Authorization"] != nil {

				tokenStr := r.Header.Get("Authorization")
				verifyResult, err := VerifyJWT(tokenStr, sign)

				if verifyResult.Result {
					next(w, r)
				} else {
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("You're Unauthorized due to invalid token"))
				}

			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Need authorization header with token !!!"))
			}

		})

	}

}

type VerifyResult struct {
	Result bool
	Token  *jwt.Token
}

/*
Метод проверки токена если true ,то валидный иначе не валидный

Использование: передаем строку jwt токена и подпись сервиса который мы используем
*/
func VerifyJWT(tokenStr string, sign []byte) (*VerifyResult, error) {
	t, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return sign, nil

	})

	if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return &VerifyResult{
			Result: true,
			Token:  t,
		}, nil
	} else {
		return &VerifyResult{
			Result: false,
			Token:  t,
		}, err

	}

}
