package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/m-a-r-a-t/go-jwt-auth/internal/models"
)

type AuthRepo struct {
	Db *sqlx.DB
}

func (ar *AuthRepo) GetUserByEmail(email string) (*models.UserData, error) {

	var user models.UserData
	err := ar.Db.Get(&user,
		`SELECT * FROM "User" WHERE email=$1`,
		email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ar *AuthRepo) GetUserRefreshTokenById(user_id []uint8) (string, error) {

	var refreshToken string
	err := ar.Db.Get(&refreshToken,
		`SELECT refresh_token FROM "UserRefreshTokens" WHERE user_id=$1`,
		user_id,
	)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (ar *AuthRepo) CreateUser(userData *models.UserData) (bool, error) {
	// ! Сделать это все в транзакции

	_, err := ar.Db.NamedExec(`INSERT INTO "User" (email) VALUES (:email)`, userData)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (ar *AuthRepo) CreateOrUpdateUserRefreshToken(refreshToken string, user_id []uint8) bool {
	// ! Сделать это все в транзакции
	result := ar.Db.MustExec(
		`INSERT INTO "UserRefreshTokens" (user_id, refresh_token) 
	VALUES ($1,$2)
	ON CONFLICT (user_id) DO UPDATE 
		SET refresh_token = $2`,
		user_id,
		refreshToken,
	)
	fmt.Println(result.RowsAffected())
	fmt.Println("Result of inserting resfresh token", result)

	return true
}
