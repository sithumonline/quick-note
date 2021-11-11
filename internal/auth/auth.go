package auth

import (
	"errors"
	"net/http"
	"os"

	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/sithumonline/quick-note/transact/user"

	"github.com/sithumonline/quick-note/models"
)

type Auth struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func (a *Auth) TokenValidation(r *http.Request, verification bool) (bool, int, error) {
	jwtKey := []byte(os.Getenv("SECRET_KEY"))
	tknStr := r.Header.Get("Authorization")

	if len(tknStr) == 0 {
		return false, http.StatusBadRequest, errors.New("token not found")
	}

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, http.StatusUnauthorized, err
		}

		return false, http.StatusBadRequest, err
	}

	if !tkn.Valid {
		err := errors.New("token is invalid")
		return false, http.StatusUnauthorized, err
	}

	userRepo := user.NewUserRepo(a.db)

	p, err := userRepo.GetPasswordByEmail(&models.Credentials{Email: claims.Email, Password: claims.Password}, verification)

	if err != nil {
		return false, http.StatusUnauthorized, err
	}

	if p != claims.Password {
		return false, http.StatusUnauthorized, errors.New("hash and password is not matching")
	}

	return true, http.StatusOK, nil
}
