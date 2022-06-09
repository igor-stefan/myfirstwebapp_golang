package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/render"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var appConfig config.AppConfig
var mySession *scs.SessionManager
var pathToTemplates = "./../../templates"

func TestMain(m *testing.M) {
	gob.Register(models.Reserva{})
	//mudar para true quando estiver em producao
	appConfig.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	mySession = scs.New()
	mySession.Lifetime = 24 * time.Hour
	mySession.Cookie.Persist = true
	mySession.Cookie.SameSite = http.SameSiteLaxMode
	mySession.Cookie.Secure = appConfig.InProduction
	appConfig.Session = mySession

	render.SetConfig(&appConfig)

	// tc é um mapa que armazena todos os templates html;
	// erro armazena possiveis erros que possam ocorrer no processamento dos templates html
	tc, erro := CreateTestTemplateCache()
	if erro != nil {
		log.Fatal("Nao foi possivel carregar os templates", erro)
	}

	// depois de carregados os templates, eles sao armazenados na variavel appConfig
	appConfig.TemplateCache = tc
	appConfig.UseCache = true //definido como false pois esta em desenvolvimento

	SetRepo(testNewRepo(&appConfig)) //passa as configs para o pkg handlers
	render.SetConfig(&appConfig)     //passa as configs para o pkg render

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/home", Repo.Home)
	mux.Get("/catalogo", Repo.Catalogo)
	mux.Post("/catalogo", Repo.PostCatalogo)
	mux.Post("/catalogo-json", Repo.CatalogoJson)
	mux.Get("/nbagame", Repo.NbaGame)
	mux.Get("/info", Repo.Info)
	mux.Get("/sb", Repo.Sb)
	mux.Get("/jancb", Repo.JanelaCopacabana)
	mux.Get("/reserva", Repo.Reserva)
	mux.Post("/reserva", Repo.PostReserva)
	mux.Get("/resumo-reserva", Repo.ResumoReserva)
	//arquivos é uma variavel que guarda os arquivos estaticos e os fornece para a pag
	arquivos := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", arquivos))

	return mux
}

// NoSurf adiciona protecao CSRF para todas as requisicoes POST
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad carrega e salva a sessão em cada requisicao
func SessionLoad(next http.Handler) http.Handler {
	return mySession.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pags, erro := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if erro != nil {
		return myCache, erro
	}

	if !(appConfig.InProduction) {
		fmt.Print("Encontrados esses arquivos para template => ")
	}
	for j, pag := range pags {
		name := filepath.Base(pag)
		if !appConfig.InProduction {
			fmt.Print(name, " | ")
			if j == len(pags)-1 {
				fmt.Println()
			}
		}
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
