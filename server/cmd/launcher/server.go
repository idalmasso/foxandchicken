package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idalmasso/foxandchicken/server/game"
	"github.com/idalmasso/foxandchicken/server/login"
)

func main() {
	gameInstance := game.NewInstance()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	var loginRouter login.UserAuthRouter

	loginRouter.LoginUser = gameInstance.AddPlayer
	r.Post("/login", loginRouter.ManageRequest)
	http.ListenAndServe(":3000", r)
}
