package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idalmasso/foxandchicken/server/game"
	"github.com/idalmasso/foxandchicken/server/gameserver"
)

func main() {
	gameInstance := game.NewInstance()
	r := chi.NewRouter()
	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	var webServer gameserver.GameServer
	webServer.Instance = gameInstance
	go gameInstance.GameInstanceRun()

	FileServer(r)
	r.Get("/api/ws", webServer.ManageRequest)
	panic(server.ListenAndServe())
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "../../../client/dist"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
