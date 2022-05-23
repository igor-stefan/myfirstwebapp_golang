package main

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig
	ret := routes(app)
	switch ret.(type) {
	case *chi.Mux:
		//ok
	default:
		t.Errorf("o tipo retornado é %T, mas é esperado *chi.Mux", ret)
	}
}
