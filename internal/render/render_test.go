package render

import (
	"net/http"
	"testing"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

func TestAdicionarDadosDefault(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	testSession.Put(r.Context(), "flash", "123")
	result := AdicionarDadosDefault(&td, r)
	if result.Flash != "123" {
		t.Error("Teste falhou, valor '123' nao encontrado")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/texto", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = testSession.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	testApp.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	pagsParaTeste := []string{"catalogo", "home", "info", "janela-copacabana", "nbagame", "reserva", "sao-bernardo", "nao-existente"}

	for _, nome := range pagsParaTeste {
		err = RenderTemplate(&ww, r, nome+".page.html", &models.TemplateData{})
		if err != nil { // se houve um erro no teste
			if nome != pagsParaTeste[len(pagsParaTeste)-1] { // e nao foi no ultimo teste
				t.Error("Erro ao renderizar os templates no browser")
			}
			//somente é especificado o último caso pois espera-se que ele retorne um erro
		}
	}

}

func TestSetConfigForRenderPkg(t *testing.T) {
	SetConfigForRenderPkg(&testApp)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("nao foi possivel criar template cache")
	}
}
