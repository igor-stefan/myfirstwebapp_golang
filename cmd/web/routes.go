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
		r.Get("/dashboard", handlers.Repo.AdminDashboard)      // pag inicial do dashboard admin
		r.Get("/livros", handlers.Repo.AdminPagLivros)         // pag mostrando todos os livros disponiveis
		r.Get("/reservas/all", handlers.Repo.AdminReservasAll) // pag mostrando todas as reservas
		r.Get("/reservas/new", handlers.Repo.AdminReservasNew) // pag mostrando apenas reservas nao processadas

		r.Get("/calendario", handlers.Repo.AdminCalendario)      // pag calendario
		r.Post("/calendario", handlers.Repo.AdminPostCalendario) // salvar alteracoes na pag calendario (add blocks)

		r.Get("/reservas/{src}/{id}/show", handlers.Repo.AdminShowReserva) // handler da pag de uma reserva especifica
		r.Post("/reservas/{src}/{id}", handlers.Repo.AdminPostShowReserva) // handler post de uma reserva especifica

		r.Get("/processarReserva/{src}/{id}/do", handlers.Repo.AdminProcessarReserva) // handler requis. marcar como proc.
		r.Get("/deletarReserva/{src}/{id}/do", handlers.Repo.AdminDeletarReserva)     // handler requis. deletar reserva especifica

	})

	return mux
}
