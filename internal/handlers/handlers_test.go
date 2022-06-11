package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
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
	{"post-reserva", "/reserva", "POST", http.StatusOK},
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

func TestRepository_Reserva(t *testing.T) {
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
		t.Errorf("Handler da pag Reserva retornou código de resposta errado, retornou %d, esperado %d", rr.Code, http.StatusOK)
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

func TestRepository_PostReserva(t *testing.T) {
	// é preciso ter alguns dados referentes à reserva já armazenados na session
	layout := "02/01/2006" // para conversao da string em time.Time
	di := "01-01-2099"
	df := "01-01-2100"
	di_t, _ := helpers.ConvStr2Time(layout, di)
	df_t, _ := helpers.ConvStr2Time(layout, df)
	DadosReserva := models.Reserva{ //cria um modelo de Reserca para ser colocado na session
		LivroID:    1000,
		DataInicio: di_t,
		DataFinal:  df_t,
		Livro: models.Livro{
			ID:        1,
			NomeLivro: "Uma Janela em Copacabana",
		},
	}

	//é preciso também construir o body do form
	// a sequencia abaixo junta a string atual seu novo valor e separa pelo caractere '&'
	reqBody := "data_inicio=01-01-2099"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "data_final=01-01-2100")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "nome=Jaylen")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "sobrenome=Brown")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jb@celtics.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "livro_id=100")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "obs=O Homem Mau")

	//io.Reader allows you to read data from something that implements the io.Reader interface into a slice of bytes
	req, _ := http.NewRequest("POST", "/reserva", strings.NewReader(reqBody)) // cria uma request
	// colocar a variavel reserva na sessão da request -> usar context
	ctx := getCtx(req)                                   // ctx que pode ser adicionado na request
	req = req.WithContext(ctx)                           // returns a shallow copy of r with its context changed to ctx
	mySession.Put(ctx, "infoReservaAtual", DadosReserva) // adiciona os dados necessarios na session

	// setar Header da request (excellent practice)
	//indica ao Browser qual o tipo da request que está chegando
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReserva)
	handler.ServeHTTP(rr, req)

	// qual é o tipo de retorno da funcao se tudo ocorrer bem? status see other, portanto, checá-lo
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReserva retornou código %d, esperado é %d", rr.Code, http.StatusSeeOther)
	}

	////////// teste para post body ausente ////////////////

	req, _ = http.NewRequest("POST", "/reserva", nil) // cria uma request
	ctx = getCtx(req)                                 // ctx que pode ser adicionado na request
	req = req.WithContext(ctx)                        // returns a shallow copy of r with its context changed to ctx
	// setar Header da request (excellent practice)
	//indica ao Browser qual o tipo da request que está chegando
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReserva)
	handler.ServeHTTP(rr, req)

	// qual é o tipo de retorno da funcao se tudo ocorrer bem? status see other, portanto, checá-lo
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReserva retornou código %d, esperado é %d quando ocorre body ausente", rr.Code, http.StatusTemporaryRedirect)
	}

	////////// teste session dados ausentes////////////////

	req, _ = http.NewRequest("POST", "/reserva", strings.NewReader(reqBody)) // cria uma request
	ctx = getCtx(req)                                                        // ctx que pode ser adicionado na request
	req = req.WithContext(ctx)                                               // returns a shallow copy of r with its context changed to ctx
	// setar Header da request (excellent practice)
	//indica ao Browser qual o tipo da request que está chegando
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReserva)
	handler.ServeHTTP(rr, req)

	// qual é o tipo de retorno da funcao se tudo ocorrer bem? status see other, portanto, checá-lo
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReserva retornou código %d, esperado é %d para session com dados ausentes", rr.Code, http.StatusTemporaryRedirect)
	}

	////////// teste form invalido ////////////////

	//é preciso reconstruir o body do form e invalidar algum dado, no caso, o email
	reqBody = *new(string)
	reqBody = "data_inicio=01-01-2099"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "data_final=01-01-2100")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "nome=Jaylen")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "sobrenome=Brown")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jb@celtics@.@com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "livro_id=100")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "obs=O Homem Mau")

	req, _ = http.NewRequest("POST", "/reserva", strings.NewReader(reqBody)) // cria uma request
	ctx = getCtx(req)                                                        // ctx que pode ser adicionado na request
	req = req.WithContext(ctx)                                               // returns a shallow copy of r with its context changed to ctx
	mySession.Put(ctx, "infoReservaAtual", DadosReserva)                     // adiciona os dados necessarios na session
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReserva)
	handler.ServeHTTP(rr, req)

	// qual é o tipo de retorno da funcao se tudo ocorrer bem? status see other, portanto, checá-lo
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReserva retornou código %d, esperado é %d", rr.Code, http.StatusSeeOther)
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
