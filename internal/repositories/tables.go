package repositories

import "github.com/jmoiron/sqlx"

type User struct {
	Id    int
	Email string
}

var userShema = `
CREATE TABLE IF NOT EXISTS "User" (
    id UUID DEFAULT gen_random_uuid () unique  not null primary key,
    email text unique not null 
);`

func CreateUserTable(db *sqlx.DB) {
	db.MustExec(userShema)
}

var usersRefreshTokensSchema = `
CREATE TABLE IF NOT EXISTS "UserRefreshTokens" (
	user_id UUID unique not null REFERENCES "User" (id),
	refresh_token text unique not null 
);
`

func CreateUsersRefreshTokensSchemaTable(db *sqlx.DB) {
	db.MustExec(usersRefreshTokensSchema)
}