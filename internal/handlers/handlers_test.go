package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"nbagame", "/nbagame", "GET", http.StatusOK},
	{"info", "/info", "GET", http.StatusOK},
	{"reserva", "/reserva", "GET", http.StatusOK},
	{"sb", "/sb", "GET", http.StatusOK},
	{"jancb", "/jancb", "GET", http.StatusOK},
	{"catalogo", "/catalogo", "GET", http.StatusOK},
	// {"post-catalogo", "/catalogo", "POST", []postData{
	// 	{key: "inicio", value: "01-01-2020"},
	// 	{key: "end", value: "01-05-2020"},
	// }, http.StatusOK},
	// {"post-catalogo-json", "/catalogo-json", "POST", []postData{
	// 	{key: "inicio", value: "01-01-2020"},
	// 	{key: "end", value: "01-05-2020"},
	// }, http.StatusOK},
	// {"post-reserva", "/reserva", "POST", []postData{
	// 	{key: "nome", value: "Jimmy"},
	// 	{key: "sobrenome", value: "Butler"},
	// 	{key: "email", value: "jbut@nba.com"},
	// 	{key: "phone", value: "999999999"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() //fecha o servidor de testes quando a função termina

	for _, test := range theTests {
		// if test.method == "GET" {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("para %s, esperado código %d, porém foi recebido %d", test.name, test.expectedStatusCode, resp.StatusCode)
		}
		// } else { //é POST
		// 	values := url.Values{}
		// 	for _, j := range test.params {
		// 		values.Add(j.key, j.value)
		// 	}
		// 	resp, err := ts.Client().PostForm(ts.URL+test.url, values)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}

		// 	if resp.StatusCode != test.expectedStatusCode {
		// 		t.Errorf("para %s, esperado código %d, porém foi recebido %d", test.name, test.expectedStatusCode, resp.StatusCode)
		// 	}
		// }
	}

}

func TestRepositoryReservation(t *testing.T) {
	DadosReserva := models.Reserva{ //cria um modelo de Reserca para ser colocado na session
		LivroID: 1,
		Livro: models.Livro{
			ID:        1,
			NomeLivro: "Uma Janela em Copacabana",
		},
	}
	req, _ := http.NewRequest("GET", "/reserva", nil) // cria uma request
	// colocar a variavel reserva na sessão da request -> usar context
	ctx := getCtx(req)         // ctx que pode ser adicionado na request
	req = req.WithContext(ctx) // returns a shallow copy of r with its context changed to ctx

	rr := httptest.NewRecorder()

	mySession.Put(ctx, "infoReservaAtual", DadosReserva) // OBS: deve ser a mesma chave presente na funcao Reserva
	// no caso 'infoReservaAtual'

	// chamar a funcao Reserva com o método GET
	// isso irá ativar o handler responsavel
	// porém nao é possível chamá-la diretamente, é necessário torná-lo uma HandlerFunction
	handler := http.HandlerFunc(Repo.Reserva)

	handler.ServeHTTP(rr, req) // chama a funcao desejada

	if rr.Code != http.StatusOK {
		t.Errorf("Handler da Reserva retornou código de resposta errado, retornou %d, esperado %d", rr.Code, http.StatusOK)
	}

	// teste em que nao é possível resgatar as infos da reserva atual da sessão
	req, _ = http.NewRequest("GET", "/reserva", nil) // reset na requisicao
	ctx = getCtx(req)                                // garante acesso à sessão para a nova request
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("o código retornado foi %d, o código esperado é %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := mySession.Load(req.Context(), req.Header.Get("X-Session"))
	// Header X-Session é necessário para que se possa ler ou escrever na session
	if err != nil {
		log.Println(err)
	}
	return ctx
}
