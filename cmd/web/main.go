package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/handlers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/render"
)

const Porta = ":8080"

// appConfig armazena as configs do app;
// os elementos que compoem as configs estao no pkg config;
// inclui templates, se deve ser usado o cache, etc.
var appConfig config.AppConfig
var mySession *scs.SessionManager

func main() {

	gob.Register(models.Reserva{})
	//mudar para true quando estiver em producao
	appConfig.InProduction = false

	mySession = scs.New()
	mySession.Lifetime = 24 * time.Hour
	mySession.Cookie.Persist = true
	mySession.Cookie.SameSite = http.SameSiteLaxMode
	mySession.Cookie.Secure = appConfig.InProduction
	appConfig.Session = mySession

	render.SetConfigForRenderPkg(&appConfig)

	// tc é um mapa que armazena todos os templates html;
	// erro armazena possiveis erros que possam ocorrer no processamento dos templates html
	tc, erro := render.CreateTemplateCache()
	if erro != nil {
		log.Fatal("Nao foi possivel carregar os templates", erro)
	}

	// depois de carregados os templates, eles sao armazenados na variavel appConfig
	appConfig.TemplateCache = tc
	appConfig.UseCache = false //definido como false pois esta em desenvolvimento

	handlers.SetHandlersRepo(handlers.NewHandlersRepo(&appConfig)) //passa as configs para o pkg handlers
	render.SetConfigForRenderPkg(&appConfig)                       //passa as configs para o pkg render

	if !appConfig.InProduction {
		fmt.Printf("Iniciando o app na porta %s\n", Porta) //faz um log do que está ocorrendo
	}

	srv := &http.Server{ //cria uma variavel do tipo server e atribui alguns valores
		Addr:    Porta,
		Handler: routes(&appConfig),
	}

	erro = srv.ListenAndServe()                           //inicia o server
	log.Fatal("Deu problema na execução do server", erro) //anuncia um possivel erro e encerra o programa
}
