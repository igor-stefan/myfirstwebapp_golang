package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch h.(type) {
	case http.Handler:
		//passou no teste
	default:
		t.Errorf("o tipo é %T, quando na verdade é esperado http.Handler", h)
	}
}

func TestWriteToConsole(t *testing.T) {
	var myH myHandler
	h := WriteToConsole(&myH)

	switch h.(type) {
	case http.Handler: //ok
	default:
		t.Errorf("o tipo é %T, quando na verdade é esperado http.Handler", h)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)
	switch h.(type) {
	case http.Handler:
		//ok
	default:
		t.Errorf("o tipo é %T, porém o esperado é http.Handler", h)

	}
}
