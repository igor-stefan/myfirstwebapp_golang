package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/igor-stefan/myfirstwebapp_golang/pkg/config"
	"github.com/igor-stefan/myfirstwebapp_golang/pkg/models"
	"github.com/igor-stefan/myfirstwebapp_golang/pkg/render"
)

// Repo é a variável que armazena repositório usado pelos handlers;
// É atualizada toda vez que SetHandlersRepo é executada
var Repo *Repository

// Repository é a estrutura de repositorio para os handlers;
// inclui as configuracoes do app, podendo ter outras
type Repository struct {
	App *config.AppConfig
}

// NewHandlersRepo retorna uma struct do tipo Repository toda vez que é executada
func NewHandlersRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// SetHandlersRepo seta o repositorio para os handlers
func SetHandlersRepo(r *Repository) {
	Repo = r
}

// Home é o handler da pagina /Home ou /
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(fmt.Sprintf("Método POST utilizado | Data de inicio é: %s | Data final é %s", inicio, final)))
}

type RespostaJson struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) CatalogoJson(w http.ResponseWriter, r *http.Request) {
	resp := RespostaJson{
		Ok:      true,
		Message: "Disponivel",
	}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Info(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "ip_remoto", remoteIP)
	render.RenderTemplate(w, r, "info.page.html", &models.TemplateData{})
}

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

func (m *Repository) Sb(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "sao-bernardo.page.html", &models.TemplateData{})
}

func (m *Repository) JanelaCopacabana(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "janela-copacabana.page.html", &models.TemplateData{})
}

func (m *Repository) Reserva(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reserva.page.html", &models.TemplateData{})
}
