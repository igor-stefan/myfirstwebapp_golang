package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi"
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

//PostCatalogo lida com as requisiçoes post na pag catalogo
func (m *Repository) PostCatalogo(w http.ResponseWriter, r *http.Request) {
	inicio := r.Form.Get("data_inicio")
	final := r.Form.Get("data_final")

	layout := "02-01-2006"
	dataInicio, err := helpers.ConvStr2Time(layout, inicio)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, final)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	livros, err := m.DB.SearchAvailabilityForAllRooms(dataInicio, dataFinal)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, i := range livros {
		m.App.InfoLog.Println("Livro:", i.ID, i.NomeLivro)
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
type RespostaJson struct {
	Ok         bool   `json:"ok"`
	Message    string `json:"message"`
	LivroID    string `json:"livroID"`
	DataInicio string `json:"dataInicio"`
	DataFinal  string `json:"dataFinal"`
}

// CatalogoJson escreve no ResponseWriter especificado um Json referente às informações da pag Catalogo
func (m *Repository) CatalogoJson(w http.ResponseWriter, r *http.Request) {
	di := r.Form.Get("data_inicio") // resgata os dados do formulario preenchidos pelo usuario
	df := r.Form.Get("data_final")
	livroID, err := strconv.Atoi(r.Form.Get("id_livro")) //processa os dados fornecidos pelo usuario, converte string -> int
	if err != nil {                                      // checa erro
		helpers.ServerError(w, err)
		return
	}

	layout := "02/01/2006"
	dataInicio, err := helpers.ConvStr2Time(layout, di) // converte de string para time, pois é o tipo usado na query
	if err != nil {                                     // checa erro
		helpers.ServerError(w, err)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// query perguntando quantas restricoes existem para determinado periodo
	disponivel, err := m.DB.SearchAvailabilityByDatesByRoomID(dataInicio, dataFinal, livroID)
	if err != nil { // checa erro
		helpers.ServerError(w, err)
		return
	}
	var msg string
	if disponivel {
		msg = "Livro Disponível!"
	} else {
		msg = "Livro Indisponível"
	}
	resp := RespostaJson{
		Ok:         disponivel,
		Message:    msg,
		LivroID:    strconv.Itoa(livroID),
		DataInicio: di,
		DataFinal:  df,
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
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
		m.App.Session.Put(r.Context(), "erro", "não foi possivel adquirir as informacoes da reserva atual da sessão")
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
	dadosAtualReserva, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) //resgata os dados da atual reserva armazenados na sessao
	if !ok {
		m.App.ErrorLog.Println("Não foi possivel encontrar os dados da sessão --> Usuário redirecionado")
		m.App.Session.Put(r.Context(), "error", "Não foi possível obter os dados da Página")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm() // analisa o formulário preenchido
	if err != nil {
		helpers.ServerError(w, err) // em caso de erro, mostra server error
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
		dados["formPagReserva"] = dadosAtualReserva // armazena os dados da pag em uma variavel
		render.Template(w, r, "reserva.page.html", &models.TemplateData{
			Form: form,
			Data: dados,
		}) // renderiza o template especificado em caso de erro com os dados do usuario salvos
		return
	}

	newReservaID, err := m.DB.InsertReserva(dadosAtualReserva) // caso validado, insere nova reserva no db
	if err != nil {
		helpers.ServerError(w, err) // caso ocorra erro no processo, server error
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
		helpers.ServerError(w, err)
	}
	m.App.Session.Put(r.Context(), "infoReservaAtual", dadosAtualReserva) // armazena na session os dados do formulario preenchido
	http.Redirect(w, r, "/resumo-reserva", http.StatusSeeOther)           //redireciona para a tabela de reservas

}

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

	render.Template(w, r, "resumo-reserva.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      dadosPEnviarPPag,
	})
}

func (m *Repository) LivroSelecionado(w http.ResponseWriter, r *http.Request) {
	IDLivro, err := strconv.Atoi(chi.URLParam(r, "id")) // armazena em IDLivro o id presente na url
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	NomeLivro := chi.URLParam(r, "nome_livro")    // armazena em NomeLivro o nome do livro presente na url
	NomeLivro, err = url.QueryUnescape(NomeLivro) //caso a string tenha os caracteres codificados para url, fazer a decodificacao
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.InfoLog.Println("Id livro sel.=", IDLivro, "Nome do livro sel. =", NomeLivro) //somente um log
	res, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva)      // resgato as informacoes da reserva atual armazenadas na session
	if !ok {                                                                            // verifica se a conversao para o tipo especificado deu certo
		helpers.ServerError(w, err)
		return
	}
	res.LivroID = IDLivro // coloca na variavel que contem os dados da sessao atual o valor de IDLivro
	res.Livro = models.Livro{
		ID:        IDLivro,
		NomeLivro: NomeLivro,
	}

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
		helpers.ServerError(w, err)
		return
	}
	di := r.URL.Query().Get("di") //resgato valores da url
	df := r.URL.Query().Get("df") //resgato valores da url

	layout := "02/01/2006"                              //layout para conversao
	dataInicio, err := helpers.ConvStr2Time(layout, di) // converte de string para time.Time pois é o utilizado na struct
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//busca nome do livro no db
	livroNome, err := m.DB.GetLivroByID(livroID)
	if err != nil {
		helpers.ServerError(w, err)
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
