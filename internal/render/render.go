package render

import (
	"bytes"
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
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache { // verificacao se esta em modo desenvolvimento
		tc = app.TemplateCache //se nao estiver, utiliza os templates encontrados no inicio da aplicacao
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(ok)
	}

	td = AdicionarDadosDefault(td, r) //adiciona dados necessarios para seguranca como tokens

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)
	_, erro := buf.WriteTo(w)
	if erro != nil {
		fmt.Println("Erro ao escrever o template no browser", erro)
	}
}

// CreateTemplateCache procura e reune os templates dos diretórios da aplicacao;
// retorna um mapa contendo todos eles (para funcionar como se fosse memória cache)
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pags, erro := filepath.Glob("./templates/*.page.html")
	if erro != nil {
		return myCache, erro
	}

	if !(app.InProduction) {
		fmt.Print("Encontrados esses arquivos para template => ")
	}
	for j, pag := range pags {
		name := filepath.Base(pag)
		if !app.InProduction {
			fmt.Print(name, " | ")
			if j == len(pags)-1 {
				fmt.Println()
			}
		}
		ts, erro := template.New(name).Funcs(functions).ParseFiles(pag)
		if erro != nil {
			return myCache, erro
		}
		encontrados, erro := filepath.Glob("./templates/*.layout.html")
		if erro != nil {
			return myCache, erro
		}
		if len(encontrados) > 0 {
			ts, erro = ts.ParseGlob("./templates/*.layout.html")
			if erro != nil {
				return myCache, erro
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
