package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

// app armazena as configs gerais da aplicacao
var app *config.AppConfig
var pathToTemplates = "./templates"

//Seta as configurações da aplicacao para uma variavel que pode ser utilizada neste pkg render
func SetConfigForRenderPkg(a *config.AppConfig) {
	app = a
}

func AdicionarDadosDefault(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renderiza um template especificado no argumento 'tmpl' em um browser usando o ResponseWriter indicado
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache { // verificacao se esta em modo desenvolvimento
		tc = app.TemplateCache //se nao estiver, utiliza os templates encontrados no inicio da aplicacao
		fmt.Println("Foi usado o Cache")
	} else {
		tc, _ = CreateTemplateCache()
		fmt.Println("Não foi usado o Cache, houve escaneamento")
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println("template especificado não encontrado")
		return errors.New("não foi possível encontrar o arquivo no cache")
	}

	td = AdicionarDadosDefault(td, r) //adiciona dados necessarios para seguranca como tokens e outros dados que devem aparecer na pág

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)
	_, erro := buf.WriteTo(w)
	if erro != nil {
		fmt.Println("Erro ao escrever o template no browser", erro)
		return erro
	}
	return nil
}

// CreateTemplateCache procura e reune os templates dos diretórios da aplicacao;
// retorna um mapa contendo todos eles (para funcionar como se fosse memória cache)
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pags, erro := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if erro != nil {
		return myCache, erro
	}

	for _, pag := range pags {

		name := filepath.Base(pag)

		ts, erro := template.New(name).Funcs(functions).ParseFiles(pag)
		if erro != nil {
			return myCache, erro
		}
		encontrados, erro := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if erro != nil {
			return myCache, erro
		}
		if len(encontrados) > 0 {
			ts, erro = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if erro != nil {
				return myCache, erro
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
