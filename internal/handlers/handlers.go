package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	myDriver "github.com/igor-stefan/myfirstwebapp_golang/internal/driver"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/forms"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/render"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/repository"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/repository/dbrepo"
)

// Repo é a variável que armazena repositório usado pelos handlers;
// É atualizada toda vez que SetHandlersRepo é executada
var Repo *Repository

// Repository é a estrutura de repositorio para os handlers;
// inclui as configuracoes do app, podendo ter outras
type Repository struct {
	App *config.AppConfig
	DB  repository.DataBaseRepo
}

// NewRepo retorna uma struct do tipo Repository toda vez que é executada
func NewRepo(a *config.AppConfig, db *myDriver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func testNewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// SetRepo seta o repositorio para os handlers
func SetRepo(r *Repository) {
	Repo = r
}

// Home é o handler da pagina /Home ou /
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// Catalogo é o handler da pag /catalogo
func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "catalogo.page.html", &models.TemplateData{})
}

// msgErroJson escreve uma mensagem de erro no responseWriter em json utilizando a estrutura respostaJson
func msgErroJson(w http.ResponseWriter, msg string, errorCode int) {
	resp := respostaJson{
		Ok:      false,
		Message: msg,
	}
	out, _ := json.MarshalIndent(resp, "", "	")
	w.WriteHeader(errorCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//PostCatalogo lida com as requisiçoes post na pag catalogo
func (m *Repository) PostCatalogo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // verificar se ocorre erro na avaliacao do form
	if err != nil {
		// não é possível avaliar o form, então retornar o json apropriado
		msgErroJson(w, "erro ao processar o form", http.StatusUnprocessableEntity)
		return
	}

	di_string := r.Form.Get("data_inicio")
	df_string := r.Form.Get("data_final")

	layout := "02-01-2006"
	dataInicio, err := helpers.ConvStr2Time(layout, di_string)
	if err != nil {
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel processar a data informada", http.StatusUnprocessableEntity)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df_string)
	if err != nil {
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel processar a data informada", http.StatusUnprocessableEntity)
		return
	}

	livros, err := m.DB.SearchAvailabilityForAllLivros(dataInicio, dataFinal) // procura no db livros disponiveis para o itvl especificado
	if err != nil {
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel conectar-se ao db", http.StatusInternalServerError)
		return
	}

	if !m.App.InProduction { // mostra os livros disponiveis encontrados do db
		for _, i := range livros {
			m.App.InfoLog.Println("Livro:", i.ID, i.NomeLivro)
		}
	}

	if len(livros) == 0 {
		m.App.InfoLog.Println("Nao há livros disponiveis para as datas selecionadas")
		m.App.Session.Put(r.Context(), "error", "Não há livros disponiveis para as datas especificadas. Tente novamente.")
		http.Redirect(w, r, "/catalogo", http.StatusSeeOther)
		return
	}

	res := models.Reserva{
		DataInicio: dataInicio,
		DataFinal:  dataFinal,
	}
	m.App.Session.Put(r.Context(), "infoReservaAtual", res)

	livrosDispEncontrados := make(map[string]interface{})
	livrosDispEncontrados["livros"] = livros
	render.Template(w, r, "escolher-livros.page.html", &models.TemplateData{
		Data: livrosDispEncontrados,
	})
}

// RespostaJson é uma estrutura que armazena parametros de uma resposta a ser dada no formato Json
type respostaJson struct {
	Ok         bool   `json:"ok"`
	Message    string `json:"message"`
	LivroID    string `json:"livroID"`
	DataInicio string `json:"dataInicio"`
	DataFinal  string `json:"dataFinal"`
}

// CatalogoJson escreve no ResponseWriter especificado um Json referente às informações da pag Catalogo
func (m *Repository) CatalogoJson(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		msgErroJson(w, "nao foi possivel processar o form ao chegar na pag", http.StatusUnprocessableEntity)
		return
	}
	di_string := r.Form.Get("data_inicio") // resgata os dados do formulario preenchidos pelo usuario
	df_string := r.Form.Get("data_final")
	livroID, err := strconv.Atoi(r.Form.Get("id_livro")) //processa os dados fornecidos pelo usuario, converte string -> int
	if err != nil {
		// helpers.ServerError(w, err)
		msgErroJson(w, "erro ao fazer conversao string -> int", http.StatusUnprocessableEntity)
		return
	}

	layout := "02/01/2006"
	dataInicio, err := helpers.ConvStr2Time(layout, di_string) // converte de string para time, pois é o tipo usado na query
	if err != nil {                                            // checa erro
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel processar a data informada", http.StatusUnprocessableEntity)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df_string)
	if err != nil {
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel processar a data informada", http.StatusUnprocessableEntity)
		return
	}

	// query perguntando quantas restricoes existem para determinado periodo
	disponivel, err := m.DB.SearchAvailabilityByDatesByLivroID(dataInicio, dataFinal, livroID)
	if err != nil { // checa erro
		// helpers.ServerError(w, err)
		msgErroJson(w, "nao foi possivel obter informacoes do db", http.StatusInternalServerError)
		return
	}
	var msg string
	if disponivel {
		msg = "Livro Disponível!"
	} else {
		msg = "Livro Indisponível"
	}
	resp := respostaJson{
		Ok:         disponivel,
		Message:    msg,
		LivroID:    strconv.Itoa(livroID),
		DataInicio: di_string,
		DataFinal:  df_string,
	}

	resposta, _ := json.MarshalIndent(resp, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resposta)

}

// Info é o handler da pag Info
func (m *Repository) Info(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "ip_remoto", remoteIP)
	render.Template(w, r, "info.page.html", &models.TemplateData{})
}

// NbaGame é o handler da pag NbaGame
func (m *Repository) NbaGame(w http.ResponseWriter, r *http.Request) {
	intMap := make(map[string]int)
	pdallas := [7]int{114, 109, 103, 111, 80, 113, 123}
	pphx := [7]int{121, 129, 94, 101, 110, 86, 90}
	for j, pto := range pdallas {
		dal := "dal" + strconv.Itoa(j+1)
		intMap[dal] = pto
	}
	for j, pto := range pphx {
		phx := "phx" + strconv.Itoa(j+1)
		intMap[phx] = pto
	}

	stringMap := make(map[string]string)
	stringMap["vencedor"] = "Dallas Mavericks"

	remoteIP := m.App.Session.GetString(r.Context(), "ip_remoto")
	stringMap["ip_remoto"] = remoteIP
	render.Template(w, r, "nbagame.page.html", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
	})
}

// Sb é o handler da pag do livro Sao Bernardo
func (m *Repository) Sb(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "sao-bernardo.page.html", &models.TemplateData{})
}

// JanelaCopacabana é o handler da pag do livro Uma Janela em Copacabana
func (m *Repository) JanelaCopacabana(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "janela-copacabana.page.html", &models.TemplateData{})
}

// Reserva é o handler da requisicao Get da pagina de reserva
func (m *Repository) Reserva(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) // resgato as informacoes da reserva atual armazenadas na session
	if !ok {                                                                       // verifica se a conversao para o tipo especificado deu certo
		m.App.ErrorLog.Println("Não foi possivel adquirir as informacoes da reserva atual da sessão")
		m.App.Session.Put(r.Context(), "error", "Não foi possivel processar as informacoes da página")
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
		return
	}
	dadosAtualReserva := make(map[string]interface{}) // variavel para armazenar os dados
	dadosAtualReserva["formPagReserva"] = res         // atribui os dados coletados à variável por meio de uma chave

	di := res.DataInicio.Format("02/01/2006") //as datas vêm em time.Time, devem ser passadas para string
	df := res.DataFinal.Format("02/01/2006")

	mystringMap := make(map[string]string) //variavel para armazenar os valores de texto (datas)
	mystringMap["data_inicio"] = di
	mystringMap["data_final"] = df

	m.App.Session.Put(r.Context(), "infoReservaAtual", res)

	render.Template(w, r, "reserva.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: mystringMap,
		Data:      dadosAtualReserva,
	})
}

// PostReserva lida com as requisiçoes post na pag catalogo; Verifica se há erro; Armazena as infos do formulario e renderiza a pag novamente
func (m *Repository) PostReserva(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // analisa o formulário preenchido
	if err != nil {
		m.App.ErrorLog.Println("Não foi possivel processar os dados do formulário -> Usuário redirecionado\n", err)
		m.App.Session.Put(r.Context(), "error", "Não foi possível processar os dados do formulário")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	dadosAtualReserva, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) //resgata os dados da atual reserva armazenados na sessao
	if !ok {
		m.App.ErrorLog.Println("Não foi possivel encontrar os dados da sessão --> Usuário redirecionado")
		m.App.ErrorLog.Println("Erro na conversão dos dados da session para o tipo models.Reserva")
		m.App.Session.Put(r.Context(), "error", "Não foi possível obter os dados da Página")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//acrescenta na variavel que possui os dados da reserva atual as infos fornecidas pelo usuario
	dadosAtualReserva.Nome = r.Form.Get("nome")
	dadosAtualReserva.Sobrenome = r.Form.Get("sobrenome")
	dadosAtualReserva.Email = r.Form.Get("email")
	dadosAtualReserva.Phone = r.Form.Get("phone")
	dadosAtualReserva.Obs = r.Form.Get("obs")

	form := forms.New(r.PostForm)                        // cria um form para ser postado na pag
	form.Required("nome", "sobrenome", "email", "phone") // chama verificacao
	form.TamMin("nome", 3)                               // verifica tamanho do campo 'nome'
	form.IsEmail("email")                                //verifica o campo email

	if !form.Valid() { // verifica se ocorreram erros no form
		dados := make(map[string]interface{})
		dados["formPagReserva"] = dadosAtualReserva // armazena os dados da pag em uma variavel para renderizá-los
		m.App.ErrorLog.Println("nao foi possivel validar o formulario, verifique os dados inseridos")
		http.Error(w, "nao foi possível processar o formulario", http.StatusSeeOther) // informa erro para a requisicao http
		render.Template(w, r, "reserva.page.html", &models.TemplateData{
			Form: form,
			Data: dados,
		}) // renderiza o template especificado em caso de erro com os dados do usuario salvos
		return
	}

	newReservaID, err := m.DB.InsertReserva(dadosAtualReserva) // caso validado, insere nova reserva no db
	if err != nil {
		m.App.ErrorLog.Println("Não foi possivel inserir nova reserva no db\n", err)
		m.App.Session.Put(r.Context(), "error", "Não foi possível realizar a reserva")
		http.Error(w, "nao foi possivel inserir a reserva no db", http.StatusTemporaryRedirect)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// helpers.ServerError(w, err) // caso ocorra erro no processo, server error
		return
	}

	restricao := models.LivroRestricao{ // cria uma variavel para armazenar dados que serao upados na tabela LivrosRestricoes do db
		DataInicio:  dadosAtualReserva.DataInicio,
		DataFinal:   dadosAtualReserva.DataFinal,
		LivroID:     dadosAtualReserva.LivroID,
		ReservaID:   newReservaID, // novo numero de reserva que retornou da funcao InsertReserva
		RestricaoID: 1,            // tipo da restricao é 1 'emprestimo_usuario'
	}
	err = m.DB.InsertLivroRestricao(restricao) // insere uma nova linha no db na tabela LivrosRestricoes
	if err != nil {                            // checa se houve erros
		m.App.ErrorLog.Println("Não foi possivel inserir nova restricao para um livro no db\n", err)
		m.App.Session.Put(r.Context(), "error", "Não foi possível realizar a reserva")
		http.Error(w, "nao foi possivel inserir a restricao do livro no db", http.StatusTemporaryRedirect)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// helpers.ServerError(w, err) // caso ocorra erro no processo, server error
		return
	}

	m.App.Session.Put(r.Context(), "infoReservaAtual", dadosAtualReserva) // armazena na session os dados do formulario preenchido
	http.Redirect(w, r, "/resumo-reserva", http.StatusSeeOther)           //redireciona para a tabela de reservas

}

// ResumoReserva é apresenta um resumo com as informacoes da Reserva
func (m *Repository) ResumoReserva(w http.ResponseWriter, r *http.Request) {
	dadosAtualReserva, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) // resgata os dados do form armazenados na sessão
	if !ok {                                                                                     // em caso de erro
		m.App.ErrorLog.Println("Não foi possivel encontrar os dados da sessão --> Usuário redirecionado")
		m.App.Session.Put(r.Context(), "error", "Não foi possível obter os dados da Página")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "infoReservaAtual") // remove os dados dos forms da sessão

	di := dadosAtualReserva.DataInicio.Format("02/01/2006") // converte de time.Time para string para ser mostrada na pag
	df := dadosAtualReserva.DataFinal.Format("02/01/2006")

	dadosPEnviarPPag := make(map[string]interface{})
	dadosPEnviarPPag["dadosAtualReserva"] = dadosAtualReserva

	stringMap := make(map[string]string) // cria string map para passar o valor das datas
	stringMap["data_inicio"] = di
	stringMap["data_final"] = df

	//envia notificacao por email para o usuario
	htmlMsg := fmt.Sprintf(`<strong>Confirmacao de Reserva</strong>
		Caro %s, <br>
		Esta é a confirmacao de sua reserva!
		Entre %s e %s`, dadosAtualReserva.Nome,
		dadosAtualReserva.DataInicio.Format("02/01/2006"),
		dadosAtualReserva.DataFinal.Format("02/01/2006"))

	msg := models.MailData{ // constroi a msg no formato esperado
		To:      dadosAtualReserva.Email,
		From:    "owner@olbookshelf.com",
		Subject: "Reserva de livro",
		Content: htmlMsg,
	}

	m.App.MailChan <- msg // envia ao channel, há um listener

	// envia msg para o administrador
	htmlMsg = fmt.Sprintf(`<h4>Reserva realizada</h4>
		Usuário: %s %s <br>
		Data Inicial: %s e Data Final: %s`,
		dadosAtualReserva.Nome, dadosAtualReserva.Sobrenome,
		dadosAtualReserva.DataInicio.Format("02/01/2006"),
		dadosAtualReserva.DataFinal.Format("02/01/2006"))

	msg = models.MailData{
		To:      "owner@olbookshelf.com",
		From:    "owner@olbookshelf.com",
		Subject: "Nova reserva realizada!",
		Content: htmlMsg,
	}
	m.App.MailChan <- msg

	render.Template(w, r, "resumo-reserva.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      dadosPEnviarPPag,
	})
}

// LivroSelecionado é o handler para uma pag intermediaria entre a escolha do livro e o formulario da Reserva
func (m *Repository) LivroSelecionado(w http.ResponseWriter, r *http.Request) {
	urlInteira := strings.Split(r.RequestURI, "/")
	id := urlInteira[2]
	IDLivro, err := strconv.Atoi(id) // armazena em IDLivro o id presente na url
	if err != nil {
		// helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "faltando parametro na url")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) // resgato as informacoes da reserva atual armazenadas na session
	if !ok {                                                                       // verifica se a conversao para o tipo especificado deu certo
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter dados da session")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
		// helpers.ServerError(w, err)
	}
	livroReq, err := m.DB.GetLivroByID(IDLivro)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter dados da requisicao")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	res.LivroID = IDLivro // coloca na variavel que contem os dados da sessao atual o valor de IDLivro
	res.Livro = livroReq

	// atualiza a sessao atual, agora possui o id do livro selecionado
	// também o tipo livro foi passado
	m.App.Session.Put(r.Context(), "infoReservaAtual", res)

	// redireciona o usuario para a pag de reserva
	http.Redirect(w, r, "/reserva", http.StatusSeeOther)
}

// ReservarLivro capta os paramatros da URL, cria a struct reserva e coloca os dados na sessão
func (m *Repository) ReservarLivro(w http.ResponseWriter, r *http.Request) {
	livroID, err := strconv.Atoi(r.URL.Query().Get("id")) //resgato os valores da url
	if err != nil {                                       // checo erros
		// helpers.ServerError(w, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	di := r.URL.Query().Get("di") //resgato valores da url
	df := r.URL.Query().Get("df") //resgato valores da url

	layout := "02/01/2006"                              //layout para conversao
	dataInicio, err := helpers.ConvStr2Time(layout, di) // converte de string para time.Time pois é o utilizado na struct
	if err != nil {
		// helpers.ServerError(w, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df)
	if err != nil {
		// helpers.ServerError(w, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//busca nome do livro no db
	livroNome, err := m.DB.GetLivroByID(livroID)
	if err != nil {
		// helpers.ServerError(w, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res := models.Reserva{ //struct para colocar os valores na sessão
		LivroID:    livroID,
		DataInicio: dataInicio,
		DataFinal:  dataFinal,
		Livro: models.Livro{
			ID:        livroID,
			NomeLivro: livroNome.NomeLivro,
		},
	}

	m.App.Session.Put(r.Context(), "infoReservaAtual", res)
	http.Redirect(w, r, "/reserva", http.StatusSeeOther)
}

// ShowLogin renderiza o formulario de login
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin trata do direcionamento do usuario após preenchimento do formulario de login
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context()) // sempre que for realizado login ou logout, renovar o token
	err := r.ParseForm()                      // analise inicial dos dados do form
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

	var email, senha string = r.PostForm.Get("email"), r.PostForm.Get("senha") // armazena os dados nas variaveis especificas
	form := forms.New(r.PostForm)
	form.Required("email", "senha") // validacao secundaria
	form.IsEmail("email")
	if !form.Valid() { // em caso de invalido
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Autenticar(email, senha) // pesquisa no db o usuario
	if err != nil {                             // caso ocorra um problema, retorna o usuario para a pag de login
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Credenciais de Login Inválidas")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	// caso tudo ok com a autenticacao
	m.App.Session.Put(r.Context(), "flash", "Logado com successo.")
	m.App.Session.Put(r.Context(), "user_id", id)
	http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusSeeOther)
}

// Logout desloga o usuário
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	m.App.Session.Destroy(r.Context())    // deleta os dados da session
	m.App.Session.RenewToken(r.Context()) // renova o token csrf
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

// AdminDashboard renderiza o dashboard de admin com suas opcoes
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

// AdminPagLivros renderiza a pag com opcoes para tratar dos livros
func (m *Repository) AdminPagLivros(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-livros.page.html", &models.TemplateData{})
}

// AdminReservasAll renderiza a pag com opcoes para tratar dos livros
func (m *Repository) AdminReservasAll(w http.ResponseWriter, r *http.Request) {
	reservas, err := m.DB.AllReservas()
	if err != nil {
		m.App.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	dadosSolicitadosDb := make(map[string]interface{})
	dadosSolicitadosDb["reservas"] = reservas
	render.Template(w, r, "admin-reservas.page.html", &models.TemplateData{
		Data: dadosSolicitadosDb,
	})
}

// AdminReservasNew mostra todas as novas reservas
func (m *Repository) AdminReservasNew(w http.ResponseWriter, r *http.Request) {
	reservas, err := m.DB.NewReservas()
	if err != nil {
		m.App.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	dadosSolicitadosDb := make(map[string]interface{})
	dadosSolicitadosDb["reservas"] = reservas
	render.Template(w, r, "admin-novas-reservas.page.html", &models.TemplateData{
		Data: dadosSolicitadosDb,
	})
}

// AdminShowReserva mostra especificações de uma única reserva
func (m *Repository) AdminShowReserva(w http.ResponseWriter, r *http.Request) {
	// pega os parametros presentes na url, encontra qual o id da reserva, query ao db
	urlDividida := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlDividida[len(urlDividida)-1])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter os dados da requisicao")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	rowReserva, err := m.DB.GetReservaById(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter os dados da requisicao")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	dadosReserva := make(map[string]interface{}, 3)
	dadosReserva["dadosReserva"] = rowReserva

	// pega a fonte da requisicao da reserva, se veio de (all ou new)
	src := urlDividida[len(urlDividida)-2]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	// renderiza a pag com os dados solicitados passados para la
	render.Template(w, r, "admin-mostra-reserva.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      dadosReserva,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReserva faz o update dos dados postados na reserva indicada
func (m *Repository) AdminPostShowReserva(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter os dados do form")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	// pega os parametros presentes na url, encontra qual o id da reserva, query ao db
	urlDividida := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(urlDividida[len(urlDividida)-1])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter os dados da requisicao")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	rowReserva, err := m.DB.GetReservaById(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel obter os dados da requisicao")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}

	// encontra a source
	src := urlDividida[len(urlDividida)-2]

	// realiza o update
	rowReserva.Nome = r.Form.Get("nome")
	rowReserva.Sobrenome = r.Form.Get("sobrenome")
	rowReserva.Phone = r.Form.Get("phone")
	rowReserva.Email = r.Form.Get("email")

	err = m.DB.UpdateReserva(rowReserva)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "nao foi possivel atualizar os dados da reserva")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}

	m.App.Session.Put(r.Context(), "flash", "Dados atualizados com sucesso!")
	http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusSeeOther)
}

// AdminCalendario renderiza a pag com opcoes para tratar do calendario
func (m *Repository) AdminCalendario(w http.ResponseWriter, r *http.Request) {
	//default dia atual
	now := time.Now()
	yParam := r.URL.Query().Get("y")
	mParam := r.URL.Query().Get("m")
	if yParam == "" || mParam == "" {
		yParam = now.Format("2006")
		mParam = now.Format("01")
	}
	mes, err := strconv.Atoi(mParam)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possível realizar conversao string p int")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	ano, err := strconv.Atoi(yParam)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possível realizar conversao string p int")
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	now = time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)

	prox := now.AddDate(0, 1, 0)
	ant := now.AddDate(0, -1, 0)

	proxMes := prox.Format("01")
	anoProxMes := prox.Format("2006")

	antMes := ant.Format("01")
	anoAntMes := ant.Format("2006")

	stringMap := make(map[string]string)
	stringMap["prox_mes"] = proxMes
	stringMap["ano_prox_mes"] = anoProxMes
	stringMap["ant_mes"] = antMes
	stringMap["ano_ant_mes"] = anoAntMes
	stringMap["atual_mes"], _ = helpers.ConvMonth2Text(mes)
	stringMap["atual_ano"] = now.Format("2006")
	stringMap["atual_mes_num"] = fmt.Sprintf("%02d", mes)

	anoAtual, mesAtual, _ := now.Date()
	localAtual := now.Location()
	firstOfMes := time.Date(anoAtual, mesAtual, 1, 0, 0, 0, 0, localAtual)
	lastOfMes := firstOfMes.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["dias_no_mes"] = lastOfMes.Day()
	intMap["mesAtual"] = int(mesAtual)

	livros, err := m.DB.AllLivros()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possível obter dados do db")
		m.App.InfoLog.Println(err)
		http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
	}
	dados := make(map[string]interface{})
	dados["livros"] = livros

	for _, x := range livros { // para cada livro
		reservasMap := make(map[string]int)                                // mapa com as restricoes por motivo de reserva
		blockMap := make(map[string]bool)                                  // mapa com as restricoes por motivo de bloqueio do admin
		for d := firstOfMes; !d.After(lastOfMes); d = d.AddDate(0, 0, 1) { // para todos os dias do mes atual, preencher com 0
			reservasMap[d.Format("02-01-2006")] = 0
			blockMap[d.Format("02-01-2006")] = false
		}

		// pegar restricoes que envolvam os dias do mes atual (a ser analisado)
		restricoes, err := m.DB.GetRestricoesForLivroByDate(x.ID, firstOfMes, lastOfMes)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "Não foi possível obter dados do db")
			m.App.InfoLog.Println(err)
			http.Redirect(w, r, "/loggedadmin/dashboard", http.StatusBadRequest)
		}

		for _, y := range restricoes { // para cada restricao retornada
			if y.ReservaID > 0 { //reserva
				for d := y.DataInicio; !d.After(y.DataFinal); d = d.AddDate(0, 0, 1) {
					reservasMap[d.Format("02-01-2006")] = y.ReservaID
				}
			} else { //block
				blockMap[y.DataInicio.Format("02-01-2006")] = true
			}
		}
		dados[fmt.Sprintf("reservas_map_%d", x.ID)] = reservasMap                   // armazena as reservas para este livro
		dados[fmt.Sprintf("block_map_%d", x.ID)] = blockMap                         // armazena os blocks para este livro
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap) // coloca na session os valores atuais de block que podem ser alterados quando houver o post
	}
	render.Template(w, r, "admin-calendario.page.html", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
		Data:      dados,
	})
}

// AdminPostCalendario lida com alterações realizadas na pagina do calendario
func (m *Repository) AdminPostCalendario(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // formulario é analisado
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
		http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
	}
	ano, err := strconv.Atoi(r.Form.Get("ano_atual")) // ano visualizado esta presente na url
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
		http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
	}
	mes, err := strconv.Atoi(r.Form.Get("mes_atual")) // mes visualizado esta presente na url
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
		http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
	}

	livros, err := m.DB.AllLivros() // acessado todos os livros
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
		http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
	}

	form := forms.New(r.PostForm) // garante acesso aos metodos do form criados no package form
	for _, x := range livros {    // para cada livro
		sessionBlockMap, ok := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]bool)
		if !ok {
			m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
			http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
		}
		for k, v := range sessionBlockMap { // para cada elemento no map da session
			if v && !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, k)) { // se havia block na session e nao há block após o post
				dataFormatoPadrão, _ := helpers.ConvStr2Time("02-01-2006", k)
				err = m.DB.DeleteBlockForLivro(x.ID, dataFormatoPadrão)
				if err != nil {
					m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
					m.App.InfoLog.Println(err)
					http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
				}
			} else if !v && form.Has(fmt.Sprintf("add_block_%d_%s", x.ID, k)) { // se não havia block na session e há block após o post
				dataFormatoPadrão, _ := helpers.ConvStr2Time("02-01-2006", k)
				err = m.DB.InsertBlockForLivro(x.ID, dataFormatoPadrão)
				if err != nil {
					m.App.Session.Put(r.Context(), "error", "Não foi possivel salvar as alterações")
					m.App.InfoLog.Println(err)
					http.Redirect(w, r, "/loggedadmin/calendario", http.StatusBadRequest)
				}
			}
		}
	}
	m.App.Session.Put(r.Context(), "flash", "Alterações salvas")
	http.Redirect(w, r, fmt.Sprintf("/loggedadmin/calendario?y=%d&m=%d", ano, mes), http.StatusSeeOther)
}

func (m *Repository) AdminProcessarReserva(w http.ResponseWriter, r *http.Request) {
	urlDividida := strings.Split(r.RequestURI, "/")
	src := urlDividida[len(urlDividida)-2]
	id, err := strconv.Atoi(urlDividida[len(urlDividida)-1])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "não foi possível obter a url solicitada")
		http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusBadRequest)
	}
	err = m.DB.UpdateProcessadaForReserva(id, 1)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Não foi possível fazer consulta ao db")
		http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusBadRequest)
	}
	m.App.Session.Put(r.Context(), "flash", "Reserva marcada como processada")
	http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeletarReserva(w http.ResponseWriter, r *http.Request) {
	urlDividida := strings.Split(r.RequestURI, "/")
	src := urlDividida[len(urlDividida)-2]
	id, err := strconv.Atoi(urlDividida[len(urlDividida)-1])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "não foi possível obter a url solicitada")
		http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusBadRequest)
	}
	err = m.DB.DeleteReserva(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "não foi possível apagar a reserva")
		http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusBadRequest)
	}
	m.App.Session.Put(r.Context(), "flash", "Reserva deletada")
	http.Redirect(w, r, fmt.Sprintf("/loggedadmin/reservas/%s", src), http.StatusSeeOther)
}
