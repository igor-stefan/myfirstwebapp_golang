package handlers

import (
	"encoding/json"
	"fmt"
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
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// CatalogoJson escreve no ResponseWriter especificado um Json referente às informações da pag Catalogo
func (m *Repository) CatalogoJson(w http.ResponseWriter, r *http.Request) {
	resp := RespostaJson{
		Ok:      true,
		Message: "Disponivel",
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
		err := fmt.Errorf("erro na conversao do valor retornado pela funcao 'Get()' da session") //cria msg de erro
		helpers.ServerError(w, err)
		return
	}
	dadosReserva := make(map[string]interface{})
	dadosReserva["formPagReserva"] = res

	di := res.DataInicio.Format("01/02/2006")
	df := res.DataFinal.Format("01/02/2006")

	mystringMap := make(map[string]string)
	mystringMap["data_inicio"] = di
	mystringMap["data_final"] = df

	render.Template(w, r, "reserva.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: mystringMap,
		Data:      dadosReserva,
	})
}

// PostReserva lida com as requisiçoes post na pag catalogo; Verifica se há erro; Armazena as infos do formulario e renderiza a pag novamente
func (m *Repository) PostReserva(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // analisa o formulário preenchido
	if err != nil {
		helpers.ServerError(w, err) // em caso de erro, mostra server error
		return
	}
	di := r.Form.Get("data_inicio") // armazena dados do form no input especificado
	df := r.Form.Get("data_final")  // armazena dados do form no input especificado

	layout := "02-01-2006"                              //layout para conversao
	dataInicio, err := helpers.ConvStr2Time(layout, di) // chama funcao de conversao e verifica erros
	if err != nil {
		helpers.ServerError(w, err) // caso erro, mostra server error
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	livroID, err := strconv.Atoi(r.Form.Get("id_livro")) // converte para int o valor do campo id_livro pois o models trata o valor como int
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	dadosFormReserva := models.Reserva{ // armazena todos os dados do form em uma variavel
		Nome:       r.Form.Get("nome"),
		Sobrenome:  r.Form.Get("sobrenome"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
		Obs:        r.Form.Get("obs"),
		DataInicio: dataInicio,
		DataFinal:  dataFinal,
		LivroID:    livroID,
	}

	form := forms.New(r.PostForm)                        // cria um form para ser postado na pag
	form.Required("nome", "sobrenome", "email", "phone") // chama verificacao
	form.TamMin("nome", 3)                               // verifica tamanho do campo 'nome'
	form.IsEmail("email")                                //verifica o campo email

	if !form.Valid() { // verifica se ocorreram erros no form
		dados := make(map[string]interface{})
		dados["formPagReserva"] = dadosFormReserva // armazena os dados da pag em uma variavel
		render.Template(w, r, "reserva.page.html", &models.TemplateData{
			Form: form,
			Data: dados,
		}) // renderiza o template especificado em caso de erro com os dados do usuario salvos
		return
	}

	newReservaID, err := m.DB.InsertReserva(dadosFormReserva) // caso validado, insere nova reserva no db
	if err != nil {
		helpers.ServerError(w, err) // caso ocorra erro no processo, server error
		return
	}

	restricao := models.LivroRestricao{ // cria uma variavel para armazenar dados que serao upados na tabela LivrosRestricoes do db
		DataInicio:  dataInicio,
		DataFinal:   dataFinal,
		LivroID:     livroID,
		ReservaID:   newReservaID, // novo numero de reserva que retornou da funcao InsertReserva
		RestricaoID: 1,            // tipo da restricao é 1 'emprestimo_usuario'
	}
	err = m.DB.InsertLivroRestricao(restricao) // insere uma nova linha no db na tabela LivrosRestricoes
	if err != nil {
		helpers.ServerError(w, err)
	}
	m.App.Session.Put(r.Context(), "formPagReserva", dadosFormReserva) // armazena na session os dados do formulario preenchido
	http.Redirect(w, r, "/resumo-reserva", http.StatusSeeOther)        //redireciona para a tabela de reservas

}

func (m *Repository) ResumoReserva(w http.ResponseWriter, r *http.Request) {
	DadosReserva, ok := m.App.Session.Get(r.Context(), "formPagReserva").(models.Reserva)
	if !ok {
		m.App.ErrorLog.Println("Não foi possivel encontrar os dados da sessão --> Usuário redirecionado")
		m.App.Session.Put(r.Context(), "error", "Não foi possível obter os dados da Página")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "formPagReserva") //remove os dados dos forms da sessão
	dados := make(map[string]interface{})
	dados["formPagReserva"] = DadosReserva
	render.Template(w, r, "resumo-reserva.page.html", &models.TemplateData{
		Data: dados,
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
	// m.App.InfoLog.Println("Nome do livro sel. =", NomeLivro) //somente um log
	res, ok := m.App.Session.Get(r.Context(), "infoReservaAtual").(models.Reserva) // resgato as informacoes da reserva atual armazenadas na session
	if !ok {                                                                       // verifica se a conversao para o tipo especificado deu certo
		helpers.ServerError(w, err)
		return
	}
	res.ID = IDLivro // coloca na variavel que contem os dados da sessao atual o valor de IDLivro
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
