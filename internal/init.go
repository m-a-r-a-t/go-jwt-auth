package internal

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/m-a-r-a-t/go-jwt-auth/internal/config"
	"github.com/m-a-r-a-t/go-jwt-auth/internal/repositories"
	"github.com/m-a-r-a-t/go-jwt-auth/internal/services"
	my_jwt "github.com/m-a-r-a-t/go-jwt-auth/pkg/jwt"
	"github.com/m-a-r-a-t/go-jwt-auth/pkg/pg_database"
	"github.com/m-a-r-a-t/go-rest-wrap/pkg/router"
)

var R router.Router

var VerifyAccessToken func(next http.HandlerFunc) http.HandlerFunc

var AuthService services.AuthService
var AuthRepo repositories.AuthRepo
var db *sqlx.DB

func init() {
	R = router.NewRouter("/api", []string{})
	VerifyAccessToken = my_jwt.CreateVerifyJWT([]byte("access"))
	s := config.Settings("internal/config/.env")

	db := pg_database.Database(s.DatabaseConf)

	repositories.CreateUserTable(db)
	repositories.CreateUsersRefreshTokensSchemaTable(db)
	AuthRepo = repositories.AuthRepo{Db: db}
	AuthService = services.AuthService{AuthRepo: &AuthRepo}

}
