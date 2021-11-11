package router

import (
	"strconv"

	"github.com/go-chi/chi"

	"github.com/sithumonline/quick-note/api/handler"
	"github.com/sithumonline/quick-note/config"
	"github.com/sithumonline/quick-note/internal/auth"
	"github.com/sithumonline/quick-note/internal/mail"
	"github.com/sithumonline/quick-note/transact/user"
)

func (o Router) UserRouter() chi.Router {
	r := chi.NewRouter()
	port, err := strconv.Atoi(config.GetEnv("mail.port"))
	if err != nil {
		panic(err)
	}
	m := mail.Mail{
		SenderAddress: config.GetEnv("mail.sender"),
		SmtpHost:      config.GetEnv("mail.host"),
		SmtpPort:      port,
		SmtpPassword:  config.GetEnv("mail.password"),
	}
	userHandler := handler.NewUserHandler(user.NewUserRepo(o.db), *auth.NewAuth(o.db), *mail.NewMail(m))

	r.Post("/signup", userHandler.Signup)
	r.Post("/signin", userHandler.Signin)
	r.Get("/welcome", userHandler.Welcome)
	r.Get("/reset/{userEmail}", userHandler.Reset)
	r.Get("/verification/{userEmail}", userHandler.Verification)
	r.Put("/{id}", userHandler.UpdateUser)

	return r
}
