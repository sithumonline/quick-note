package handler

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sithumonline/quick-note/config"
	"github.com/sithumonline/quick-note/internal/auth"
	"github.com/sithumonline/quick-note/internal/mail"
	"github.com/sithumonline/quick-note/models"
	"github.com/sithumonline/quick-note/transact/user"
)

type UserHandler struct {
	repo user.UserRepo
	auth auth.Auth
	mail mail.Mail
}

func NewUserHandler(repo user.UserRepo, auth auth.Auth, mail mail.Mail) *UserHandler {
	return &UserHandler{
		repo: repo,
		auth: auth,
		mail: mail,
	}
}

func (p *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	t := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.repo.Save(&t); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := p.repo.GetTokenWithoutCred(&models.Credentials{Email: t.Email}, false)

	p.mail.Send(mail.MailBody{
		ReceiversAddress: t.Email,
		Subject:          "[Quick Note] Email Verification",
		Body: fmt.Sprintf(
			"<a href=\"%vverify/%v/%v\" rel=\"noreferrer\" target=\"_blank\">Verify</a>",
			config.GetEnv("frontend.url"),
			token,
			b64.StdEncoding.EncodeToString([]byte(t.Email)),
		),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "successfully user created")
}

func (p *UserHandler) Signin(w http.ResponseWriter, r *http.Request) {
	t := models.Credentials{}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := p.repo.GetTokenByCred(&t)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := p.repo.GetIdByEmail(&t, true)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, token+";"+b64.StdEncoding.EncodeToString([]byte(id)))
}

func (p *UserHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)

	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "Welcome to quick-note")
}

func (p *UserHandler) Reset(w http.ResponseWriter, r *http.Request) {
	bMail, _ := b64.StdEncoding.DecodeString(chi.URLParam(r, "userEmail"))
	uMail := string(bMail)

	id, err := p.repo.GetIdByEmail(&models.Credentials{Email: uMail}, true)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := p.repo.GetTokenWithoutCred(&models.Credentials{Email: uMail}, true)
	p.mail.Send(mail.MailBody{
		ReceiversAddress: uMail,
		Subject:          "[Quick Note] Password Reset",
		Body: fmt.Sprintf(
			"<a href=\"%vreset/%v/%v\" rel=\"noreferrer\" target=\"_blank\">Reset</a>",
			config.GetEnv("frontend.url"),
			id,
			token,
		),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "token sent to email")
}

func (p *UserHandler) Verification(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, false)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	bMail, _ := b64.StdEncoding.DecodeString(chi.URLParam(r, "userEmail"))
	uMail := string(bMail)
	err = p.repo.Verification(&models.Credentials{Email: uMail})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "successfully email verified")
}

func (p *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	t := models.User{}
	ID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.repo.Update(&t, ID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "successfully updated")
}
