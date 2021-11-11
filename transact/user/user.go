package user

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/sithumonline/quick-note/models"
)

type UserRepo interface {
	Save(user *models.User) error
	GetPasswordByEmail(cred *models.Credentials, verification bool) (string, error)
	GetTokenByCred(cred *models.Credentials) (string, error)
	Migrate() error
	GetTokenWithoutCred(cred *models.Credentials, verification bool) (string, error)
	Update(user *models.User, id string) error
	Verification(cred *models.Credentials) error
	GetUserByEmail(cred *models.Credentials) (*models.User, error)
	GetIdByEmail(cred *models.Credentials, verification bool) (string, error)
}

type User struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) Save(user *models.User) error {
	passwordString, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	pws := fmt.Sprintf("%s", passwordString)
	fmt.Println(pws)
	if err != nil {
		log.Panic(err)
	}
	user.Password = pws
	user.Verification = false
	if result := u.db.Create(&user); result.Error != nil {
		log.Errorf("failed to create user: %+v: %v", u, result.Error)
		return result.Error
	}

	return nil
}

func (u *User) Update(user *models.User, id string) error {
	if user.Password != "" {
		passwordString, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			log.Panic(err)
		}
		user.Password = string(passwordString)
	}
	if result := u.db.Model(&user).Where("id = ?", id).Updates(user); result.Error != nil {
		log.Errorf("failed to update user: %+v: %v", u, result.Error)
		return result.Error
	}

	return nil
}

func (u *User) Verification(cred *models.Credentials) error {
	user, err := u.GetUserByEmail(cred)
	if err != nil {
		return err
	}
	user.Verification = true

	err = u.Update(user, user.ID.String())
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUserByEmail(cred *models.Credentials) (*models.User, error) {
	user := &models.User{}
	if result := u.db.Find(user, "email = ?", cred.Email); result.Error != nil {
		log.Errorf("failed to find user: %+v: %v", u, result.Error)
		return nil, result.Error
	}

	return user, nil
}

func (u *User) GetPasswordByEmail(cred *models.Credentials, verification bool) (string, error) {
	user, err := u.GetUserByEmail(cred)
	if err != nil {
		return "", err
	}
	if !user.Verification && verification {
		return "", errors.New(fmt.Sprintf("%v not verified", cred.Email))
	}

	return user.Password, nil
}

func (u *User) GetIdByEmail(cred *models.Credentials, verification bool) (string, error) {
	user, err := u.GetUserByEmail(cred)
	if err != nil {
		return "", err
	}
	if !user.Verification && verification {
		return "", errors.New(fmt.Sprintf("%v not verified", cred.Email))
	}

	return user.ID.String(), nil
}

func (u *User) GetTokenByCred(cred *models.Credentials) (string, error) {
	p, e := u.GetPasswordByEmail(cred, true)
	if e != nil {
		return "", e
	}
	if p == "" {
		return "", errors.New(fmt.Sprintf("password not found for %v", cred.Email))
	}
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(cred.Password))
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &models.Claims{
		Email:    cred.Email,
		Password: p,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) GetTokenWithoutCred(cred *models.Credentials, verification bool) (string, error) {
	p, e := u.GetPasswordByEmail(cred, verification)
	if e != nil {
		return "", e
	}
	if p == "" {
		return "", errors.New(fmt.Sprintf("password not found for %v", cred.Email))
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &models.Claims{
		Email:    cred.Email,
		Password: p,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) Migrate() error {
	u.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := u.db.AutoMigrate(models.User{}); err != nil {
		log.Errorf("failed to migrate user: %+v: %v", models.User{}, err)
		return err
	}

	return nil
}
