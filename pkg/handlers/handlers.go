package handlers

import (
	"net/http"

	"github.com/igor-stefan/firstwebapp_golang/pkg/config"
	"github.com/igor-stefan/firstwebapp_golang/pkg/models"
	"github.com/igor-stefan/firstwebapp_golang/pkg/render"
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
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "ip_remoto", remoteIP)
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

//About é o handler da pagina /about
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["jogo1"] = "Suns 121 x 114 Mavs"
	stringMap["jogo2"] = "Suns 129 x 109 Mavs"
	stringMap["jogo3"] = "Mavs 103 x 94 Suns"
	stringMap["jogo4"] = "Mavs 111 x 101 Suns"
	stringMap["jogo5"] = "10/05/2022"
	stringMap["jogo6"] = "12/05/2022"
	stringMap["jogo7"] = "Se necessário"

	remoteIP := m.App.Session.GetString(r.Context(), "ip_remoto")
	stringMap["ip_remoto"] = remoteIP
	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
