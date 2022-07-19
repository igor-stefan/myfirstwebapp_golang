package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"padData": FormatoDataPadrao,
}

// app armazena as configs gerais da aplicacao
var app *config.AppConfig
var pathToTemplates = "./templates"

// SetConfig seta as configurações da aplicacao para uma variavel que pode ser utilizada neste pkg render
func SetConfig(a *config.AppConfig) {
	app = a
}

// FormatoDataPadrao recebe uma instancia de tempo e retorna uma string formatada no padrao
func FormatoDataPadrao(t time.Time) string {
	return t.Format("02-01-2006")
}

// AdicionarDadosDefault acrescenta dados padrão para a request realizada, os dados sao armazenados na session
func AdicionarDadosDefault(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.Autenticado = 1
	} else { // o valor default já é 0, apenas estará aqui para indicar o comportamento
		td.Autenticado = 0
	}
	return td
}

// Template renderiza um template especificado no argumento 'tmpl' em um browser usando o ResponseWriter indicado
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	if app.UseCache { // verificacao se esta em modo desenvolvimento
		tc = app.TemplateCache //se nao estiver, utiliza os templates encontrados no inicio da aplicacao
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		app.ErrorLog.Println("template especificado não encontrado")
		return errors.New("não foi possível encontrar o arquivo no cache")
	}

	// adiciona aos dados que chegaram na funcao template outros dados padrao
	// alguns necessarios para seguranca como tokens e outros dados que devem aparecer na pág, como alertas
	td = AdicionarDadosDefault(td, r)

	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)
	_, erro := buf.WriteTo(w)
	if erro != nil {
		app.ErrorLog.Println("erro ao escrever o template no browser", erro)
		return erro
	}
	return nil
}

// CreateTemplateCache procura e reune os templates dos diretórios da aplicacao;
// Retorna um mapa contendo todos eles (para funcionar como se fosse memória cache)
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pags, erro := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if erro != nil {
		app.ErrorLog.Println("nao foi possivel encontrar os templates no path especificado")
		return myCache, erro
	}

	for _, pag := range pags {

		name := filepath.Base(pag)

		ts, erro := template.New(name).Funcs(functions).ParseFiles(pag)
		if erro != nil {
			app.ErrorLog.Println("nao foi possivel criar um template a partir do arquivo especificado, verifique-o")
			return myCache, erro
		}
		encontrados, erro := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if erro != nil {
			app.ErrorLog.Println("nao foi possivel encontrar um layout para template com o caminho especificado")
			return myCache, erro
		}
		if len(encontrados) > 0 {
			ts, erro = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if erro != nil {
				app.ErrorLog.Println("nao foi possivel criar um layout de template a partir do arquivo especificado, verifique-o")
				return myCache, erro
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
