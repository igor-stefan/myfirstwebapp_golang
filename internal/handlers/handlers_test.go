package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"nbagame", "/nbagame", "GET", []postData{}, http.StatusOK},
	{"info", "/info", "GET", []postData{}, http.StatusOK},
	{"reserva", "/reserva", "GET", []postData{}, http.StatusOK},
	{"sb", "/sb", "GET", []postData{}, http.StatusOK},
	{"jancb", "/jancb", "GET", []postData{}, http.StatusOK},
	{"catalogo", "/catalogo", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() //fecha o servidor de testes quando a função termina

	for _, test := range theTests {
		if test.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("para %s, esperado código %d, porém foi recebido %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else { //é POST

		}
	}

}
