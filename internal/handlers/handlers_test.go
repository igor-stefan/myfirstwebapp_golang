package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	myDriver "github.com/igor-stefan/myfirstwebapp_golang/internal/driver"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

var theTests = []struct { // urls que nao utilizam session
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"catalogo", "/catalogo", "GET", http.StatusOK},
	{"info", "/info", "GET", http.StatusOK},
	{"nbagame", "/nbagame", "GET", http.StatusOK},
	{"sb", "/sb", "GET", http.StatusOK},
	{"jancb", "/jancb", "GET", http.StatusOK},
	{"reserva", "/reserva", "GET", http.StatusOK},
}

func getCtx(req *http.Request) context.Context {
	ctx, err := mySession.Load(req.Context(), req.Header.Get("X-Session"))
	// Header X-Session é necessário para que se possa ler ou escrever na session
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestNewRepo(t *testing.T) {
	var ret interface{} = NewRepo(&appConfig, &myDriver.DB{}) // ret recebe o retorno da funcao NewRepo
	_, ok := ret.(*Repository)                                // verifica se a interface ret possui tipo *Repository, o valor ignorado '_' é o valor de ret
	if !ok {
		t.Errorf("o tipo retornado é %T, mas é esperado *Repository", ret)
	}
}
func TestRepository_Handlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() //fecha o servidor de testes quando a função termina

	for _, test := range theTests {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("para %s, esperado código %d, porém foi recebido %d", test.name, test.expectedStatusCode, resp.StatusCode)
		}
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
	dadosReserva := models.Reserva{ //cria um modelo de Reserca para ser colocado na session
		LivroID:    1000,
		DataInicio: di_t,
		DataFinal:  df_t,
		Livro: models.Livro{
			ID:        1,
			NomeLivro: "Uma Janela em Copacabana",
		},
	}

	// é preciso também construir o body do form
	var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=01-01-2100", "nome=Jaylen", "sobrenome=Brown",
		"email=jb@celtics.com", "phone=123456789", "livro_id=100", "obs=O Homem Mau"}
	var reqBody = *new(string) // string que recebera os params da req post
	// o loop abaixo junta a string atual com o valor seguinte sendo estes separados pelo caractere '&'
	for _, p := range postBodyParams {
		reqBody = fmt.Sprintf("%s&%s", reqBody, p)
	}

loop:
	for i := 0; ; i++ { // testar 6 casos de teste alterando os pontos necessarios para testar checagens de erro
		// io.Reader allows you to read data from something that implements the io.Reader interface into a slice of bytes
		req, _ := http.NewRequest("POST", "/reserva", strings.NewReader(reqBody[1:])) // cria uma request
		// colocar a variavel reserva na sessão da request -> usar context
		ctx := getCtx(req)                                   // ctx que pode ser adicionado na request
		req = req.WithContext(ctx)                           // returns a shallow copy of r with its context changed to ctx
		mySession.Put(ctx, "infoReservaAtual", dadosReserva) // adiciona os dados necessarios na session

		// setar Header da request (excellent practice)
		// indica ao Browser qual o tipo da request que está chegando
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		rr := httptest.NewRecorder() // responseWriter para testes

		handler := http.HandlerFunc(Repo.PostReserva)

		switch i {
		case 0: // deve passar em todas as checagens de erro
			handler.ServeHTTP(rr, req)
			// qual é o tipo de retorno da funcao se tudo ocorrer bem? status see other, portanto, checá-lo
			if rr.Code != http.StatusSeeOther {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusSeeOther, i)
			}
		case 1: // body da request ausente
			req.Body = nil
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusTemporaryRedirect {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusTemporaryRedirect, i)
			}
		case 2: // dados da session ausentes
			mySession.Remove(ctx, "infoReservaAtual")
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusTemporaryRedirect {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusTemporaryRedirect, i)
			}
		case 3: // quando nao é possivel inserir a reserva no db
			dadosReserva.LivroID = -1
			mySession.Put(ctx, "infoReservaAtual", dadosReserva)
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusTemporaryRedirect {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusTemporaryRedirect, i)
			}
		case 4: // quando nao é possivel inserir a restricao do livro no db
			dadosReserva.LivroID = -2
			mySession.Put(ctx, "infoReservaAtual", dadosReserva) // alterando ctx com um dado que invalida inserção no db
			handler.ServeHTTP(rr, req)                           // a alteracao no ctx é executada antes do handler
			if rr.Code != http.StatusTemporaryRedirect {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusTemporaryRedirect, i)
			}
		case 5: // form invalido
			postBodyParams[4] = "email=jb@celtics@.@com" // referente ao email, para alterá-lo
			reqBody = *new(string)                       // cria uma string vazia (reset dos parametros da req post)
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:])) // altera o body da request antes de servir o http
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusSeeOther {
				t.Errorf("PostReserva retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusSeeOther, i)
			}
		default:
			break loop
		}
	}
}

func TestRepository_CatalogoJson(t *testing.T) {
	var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=01-01-2100", "id_livro=1"}
	var reqBody = *new(string) // string que recebera os params da req post
	for _, p := range postBodyParams {
		reqBody = fmt.Sprintf("%s&%s", reqBody, p)
	}
loop:
	for i := 0; ; i++ {
		req, _ := http.NewRequest("POST", "/catalogo-json", strings.NewReader(reqBody[1:])) // cria request excluindo o primeiro char da string
		ctx := getCtx(req)                                                                  // pega o cxt
		req = req.WithContext(ctx)                                                          // returns a shallow copy of r with its context changed to ct
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")                 // set req header ATENCÃO AQUI!!!
		rr := httptest.NewRecorder()                                                        // response recorder
		handler := http.HandlerFunc(Repo.CatalogoJson)                                      // torna o handler uma handlerFunc
		var respostaRetornada respostaJson                                                  // deve ser retornado um json, verificá-lo

		switch i {
		case 0: // sem problema
			handler.ServeHTTP(rr, req)
			err := json.Unmarshal(rr.Body.Bytes(), &respostaRetornada)
			if err != nil {
				t.Error("nao foi possivel processar o json retornado")
			}
		case 1: // erro ao avaliar form -> body ausente
			req.Body = nil
			handler.ServeHTTP(rr, req) // inicia o teste chamando fazendo a req http
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("CatalogoJson retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
		case 2: // erro ao encontrar id livro
			var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=01-01-2100", "id_livro=x"}
			var reqBody = *new(string) // string que recebera os params da req post
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:]))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("CatalogoJson retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
		case 3: // erro ao processar data inicial da reserva
			var postBodyParams = []string{"data_inicio=invalida", "data_final=01-01-2100", "id_livro=1000"}
			var reqBody = *new(string) // string que recebera os params da req post
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:]))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("CatalogoJson retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
		case 4: //  erro ao processar data final da reserva
			var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=invalida", "id_livro=1000"}
			var reqBody = *new(string) // string que recebera os params da req post
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:]))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("CatalogoJson retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
		case 5: // erro ao procurar restricoes para determinado periodo
			var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=01-01-2100", "id_livro=-2"}
			var reqBody = *new(string) // string que recebera os params da req post
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:]))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusInternalServerError {
				t.Errorf("CatalogoJson retornou código %d, esperado é %d -> caso %d", rr.Code, http.StatusInternalServerError, i)
			}
		case 6: // mesmo do case 0, tudo ok, porém nenhum livro disponivel
			var postBodyParams = []string{"data_inicio=01-01-2099", "data_final=01-01-2100", "id_livro=-1"}
			var reqBody = *new(string) // string que recebera os params da req post
			for _, p := range postBodyParams {
				reqBody = fmt.Sprintf("%s&%s", reqBody, p)
			}
			req.Body = ioutil.NopCloser(strings.NewReader(reqBody[1:]))
			handler.ServeHTTP(rr, req)
			err := json.Unmarshal(rr.Body.Bytes(), &respostaRetornada)
			if err != nil {
				t.Error("nao foi possivel processar o json retornado")
			}
		default:
			break loop
		}
	}
}

func TestRepository_PostCatalogo(t *testing.T) {
loop:
	for i := 0; ; i++ {
		var postedData url.Values = make(url.Values)
		postedData.Add("data_inicio", "01-01-2099")
		postedData.Add("data_final", "01-01-2100")
		req, _ := http.NewRequest("POST", "/catalogo", strings.NewReader(postedData.Encode())) // cria request
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // set header
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostCatalogo)
		switch i {
		case 0: // tudo ok
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusOK, i)
			}
		case 1: // erro no form
			req.Body = nil
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
		case 2: //erro na data inicial
			postedData.Set("data_inicio", "invalid")
			req.Body = ioutil.NopCloser(strings.NewReader(postedData.Encode()))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
			postedData.Set("data_inicio", "01-01-2099")
		case 3: // erro na data_final
			postedData.Set("data_final", "invalid")
			req.Body = ioutil.NopCloser(strings.NewReader(postedData.Encode()))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnprocessableEntity {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusUnprocessableEntity, i)
			}
			postedData.Set("data_final", "01-01-2100")
		case 4: // consulta nao retorna livros
			postedData.Set("data_inicio", "01-01-3000")
			req.Body = ioutil.NopCloser(strings.NewReader(postedData.Encode()))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusSeeOther {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusSeeOther, i)
			}
		case 5: // consulta ao db retorna erro
			postedData.Set("data_inicio", "01-01-1999")
			req.Body = ioutil.NopCloser(strings.NewReader(postedData.Encode()))
			handler.ServeHTTP(rr, req)
			if rr.Code != http.StatusInternalServerError {
				t.Errorf("PostCatalogo retornou %d, esperado %d, caso %d", rr.Code, http.StatusInternalServerError, i)
			}
		default:
			break loop
		}
	}
}

func TestMsgErroJson(t *testing.T) {
	rr := httptest.NewRecorder()
	msgErroJson(rr, "mensagem de teste", http.StatusOK)
	if rr.Code != http.StatusOK {
		t.Errorf("erro no codigo de resposta da mensagem json, recebido %d, esperado %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ReservarLivro(t *testing.T) {
	urlParams := map[string]string{ // parametros padrao que validam o teste
		"id": "1",
		"di": "01-01-2099",
		"df": "01-01-2100",
	}
	testCases := map[string]urlTest{
		"tudo ok": {
			urlParams,
			http.StatusSeeOther,
		},
		"sem parametros": {
			map[string]string{},
			http.StatusBadRequest,
		},
		"id invalido": {
			map[string]string{
				"id": "invalid",
				"di": "01-01-2099",
				"df": "01-01-2100",
			},
			http.StatusBadRequest,
		},
		"di invalido": {
			map[string]string{
				"id": "1",
				"di": "invalido",
				"df": "01-01-2100",
			},
			http.StatusBadRequest,
		},
		"df invalido": {
			map[string]string{
				"id": "1000",
				"di": "01-01-2099",
				"df": "invalido",
			},
			http.StatusBadRequest,
		},
		"chamada_db_erro": {
			map[string]string{
				"id": "1000",
				"di": "01-01-2099",
				"df": "01-01-2100",
			},
			http.StatusBadRequest,
		},
	}
	for nTestCase, tp := range testCases {
		req, _ := http.NewRequest("GET", "/reservar-livro", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		q := req.URL.Query()
		for k, v := range tp.testParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.ReservarLivro)
		handler.ServeHTTP(rr, req)
		if rr.Code != tp.statusCode {
			t.Errorf("ReservarLivro retornou cod. %d, esperado %d, teste -> %s", rr.Code, tp.statusCode, nTestCase)
		}
	}
}

func TestRepository_LivroSelecionado(t *testing.T) {
	dadosReserva := models.Reserva{}
	testCases := map[string]urlTestLivro{
		"tudo_ok": {
			id:         "1",
			statusCode: http.StatusSeeOther,
		},
		"erro_no_id": {
			id:         "invalido",
			statusCode: http.StatusBadRequest,
		},
		"erro_session": {
			id:         "1",
			statusCode: http.StatusBadRequest,
		},
		"erro_id_vazio": {
			id:         "",
			statusCode: http.StatusBadRequest,
		},
		"ret_erro_db": {
			id:         "10", // id > 7 retorna erro do db
			statusCode: http.StatusBadRequest,
		},
	}

	for tc, tp := range testCases {
		q := fmt.Sprintf("/livro-selecionado/%s", tp.id)
		req, _ := http.NewRequest("GET", q, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		mySession.Put(ctx, "infoReservaAtual", dadosReserva) // adiciona os dados necessarios na session
		if tc == "erro_session" {
			mySession.Remove(ctx, "infoReservaAtual")
		}
		req.RequestURI = q
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.LivroSelecionado)
		handler.ServeHTTP(rr, req)
		if rr.Code != tp.statusCode {
			t.Errorf("LivroSelecionado retornou cod. %d, esperado %d, teste -> %s", rr.Code, tp.statusCode, tc)
		}
	}
}

func TestRepository_ResumoReserva(t *testing.T) {
	dadosReserva := models.Reserva{}
	testCases := []struct {
		nomeTeste      string
		method         string
		url            string
		hasSessionData bool
		statusCode     int
	}{
		{"tudo_ok", "GET", "/resumo-reserva", true, http.StatusOK},
		{"erro_session", "GET", "/resumo-reserva", false, http.StatusTemporaryRedirect},
	}
	for _, tp := range testCases {
		req, _ := http.NewRequest(tp.method, tp.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		if tp.hasSessionData {
			mySession.Put(ctx, "infoReservaAtual", dadosReserva)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.ResumoReserva)
		handler.ServeHTTP(rr, req)
		if rr.Code != tp.statusCode {
			t.Errorf("ResumoReserva retornou code %d, esperado %d, caso -> %s", rr.Code, tp.statusCode, tp.nomeTeste)
		}
	}
}
