package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/handlers"
)

// routes configura as rotas para cada pag da aplicacao
func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/home", handlers.Repo.Home)
	mux.Get("/catalogo", handlers.Repo.Catalogo)
	mux.Post("/catalogo", handlers.Repo.PostCatalogo)
	mux.Post("/catalogo-json", handlers.Repo.CatalogoJson)
	mux.Get("/nbagame", handlers.Repo.NbaGame)
	mux.Get("/info", handlers.Repo.Info)
	mux.Get("/sb", handlers.Repo.Sb)
	mux.Get("/jancb", handlers.Repo.JanelaCopacabana)
	mux.Get("/reserva", handlers.Repo.Reserva)
	mux.Post("/reserva", handlers.Repo.PostReserva)
	//arquivos é uma variavel que guarda os arquivos estaticos e os fornece para a pag
	arquivos := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", arquivos))

	return mux
}
