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
