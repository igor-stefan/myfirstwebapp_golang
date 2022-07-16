package main

import (
	"fmt"
	"net/http"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
	"github.com/justinas/nosurf"
)

// WriteToConsole é um middleware custom que informa que uma solicitacao chegou a pag
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !appConfig.InProduction {
			fmt.Println("Chegou na pagina")
		}
		next.ServeHTTP(w, r)
	})
}

// NoSurf adiciona protecao CSRF para todas as requisicoes POST
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad carrega e salva a sessão em cada requisicao
func SessionLoad(next http.Handler) http.Handler {
	return mySession.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.EstaAutenticado(r) {
			mySession.Put(r.Context(), "error", "O login deve ser realizado primeiro")
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
