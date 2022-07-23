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

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
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
	mux.Get("/resumo-reserva", handlers.Repo.ResumoReserva)
	mux.Get("/livro-selecionado/{id}", handlers.Repo.LivroSelecionado)
	mux.Get("/reservar-livro", handlers.Repo.ReservarLivro)

	mux.Get("/admin/login", handlers.Repo.ShowLogin)
	mux.Post("/admin/login", handlers.Repo.PostShowLogin)
	mux.Get("/admin/logout", handlers.Repo.Logout)

	//arquivos é uma variavel que guarda os arquivos estaticos e os fornece para a pag
	arquivos := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", arquivos))

	mux.Route("/loggedadmin", func(r chi.Router) {
		//r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)
		r.Get("/livros", handlers.Repo.AdminPagLivros)
		r.Get("/reservas", handlers.Repo.AdminReservas)
		r.Get("/reservas/all", handlers.Repo.AdminReservas)
		r.Get("/calendario", handlers.Repo.AdminCalendario)
		r.Post("/calendario", handlers.Repo.AdminPostCalendario)
		r.Get("/reservas/new", handlers.Repo.AdminNewReservas)
		r.Get("/reservas/{src}/{id}", handlers.Repo.AdminShowReserva)
		r.Post("/reservas/{src}/{id}", handlers.Repo.AdminPostShowReserva)
		r.Get("/processarReserva/{src}/{id}", handlers.Repo.AdminProcessarReserva)
		r.Get("/deletarReserva/{src}/{id}", handlers.Repo.AdminDeletarReserva)

	})

	return mux
}
