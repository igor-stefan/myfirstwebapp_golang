package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

// Catalogo é o handler da pag /catalogo
func (m *Repository) Catalogo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "catalogo.page.html", &models.TemplateData{})
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
	}
	w.Write([]byte(fmt.Sprintf("Método POST utilizado | Data de inicio é: %s | Data final é %s", inicio, final)))
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
	render.RenderTemplate(w, r, "info.page.html", &models.TemplateData{})
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
	render.RenderTemplate(w, r, "nbagame.page.html", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    intMap,
	})
}

// Sb é o handler da pag do livro Sao Bernardo
func (m *Repository) Sb(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "sao-bernardo.page.html", &models.TemplateData{})
}

// JanelaCopacabana é o handler da pag do livro Uma Janela em Copacabana
func (m *Repository) JanelaCopacabana(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "janela-copacabana.page.html", &models.TemplateData{})
}

// Reserva é o handler da requisicao Get da pagina de reserva
func (m *Repository) Reserva(w http.ResponseWriter, r *http.Request) {
	var vazioFormDados models.Reserva
	dados := make(map[string]interface{})
	dados["formPagReserva"] = vazioFormDados
	render.RenderTemplate(w, r, "reserva.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: dados,
	})
}

// PostReserva lida com as requisiçoes post na pag catalogo; Verifica se há erro; Armazena as infos do formulario e renderiza a pag novamente
func (m *Repository) PostReserva(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	di := r.Form.Get("data_inicio")
	df := r.Form.Get("data_final")

	layout := "2006-01-02"
	dataInicio, err := helpers.ConvStr2Time(layout, di)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	dataFinal, err := helpers.ConvStr2Time(layout, df)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	livroID, err := strconv.Atoi(r.Form.Get("id_livro"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	dadosFormReserva := models.Reserva{
		Nome:       r.Form.Get("nome"),
		Sobrenome:  r.Form.Get("sobrenome"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
		Obs:        r.Form.Get("obs"),
		DataInicio: dataInicio,
		DataFinal:  dataFinal,
		LivroID:    livroID,
	}

	form := forms.New(r.PostForm)
	form.Required("nome", "sobrenome", "email", "phone")
	form.TamMin("nome", 3)
	form.IsEmail("email")

	if !form.Valid() {
		dados := make(map[string]interface{})
		dados["formPagReserva"] = dadosFormReserva
		render.RenderTemplate(w, r, "reserva.page.html", &models.TemplateData{
			Form: form,
			Data: dados,
		})
		return
	}

	newReservaID, err := m.DB.InsertReserva(dadosFormReserva)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restricao := models.LivroRestricao{
		DataInicio:  dataInicio,
		DataFinal:   dataFinal,
		LivroID:     livroID,
		ReservaID:   newReservaID,
		RestricaoID: 1,
	}
	err = m.DB.InsertLivroRestricao(restricao)
	if err != nil {
		helpers.ServerError(w, err)
	}
	m.App.Session.Put(r.Context(), "formPagReserva", dadosFormReserva)
	http.Redirect(w, r, "/resumo-reserva", http.StatusSeeOther)

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
	render.RenderTemplate(w, r, "resumo-reserva.page.html", &models.TemplateData{
		Data: dados,
	})
}
