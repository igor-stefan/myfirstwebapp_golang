package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Error("Falhou a inicialização da aplicacao")
	}
}
