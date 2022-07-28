package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/config"
	myDriver "github.com/igor-stefan/myfirstwebapp_golang/internal/driver"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/handlers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/render"
)

const Porta = ":8080"

// appConfig armazena as configs do app;
// os elementos que compoem as configs estao no pkg config;
// inclui templates, se deve ser usado o cache, etc.
var appConfig config.AppConfig
var mySession *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, erro := run()
	if erro != nil {
		log.Fatal(erro)
	}
	defer db.SQL.Close()            // encerra a conexao com o db ao final da execucao da main
	defer close(appConfig.MailChan) // fecha o channel criado ao final da execucao da main
	listenForMail()

	if !appConfig.InProduction {
		fmt.Printf("Iniciando o app na porta %s\n", Porta) // faz um log do que está ocorrendo
	}

	srv := &http.Server{ // cria uma variavel do tipo server e atribui alguns valores
		Addr:    Porta,
		Handler: routes(&appConfig),
	}

	erro = srv.ListenAndServe()                           // inicia o server
	log.Fatal("Deu problema na execução do server", erro) // anuncia um possivel erro e encerra o programa
}

func run() (*myDriver.DB, error) {

	gob.Register(models.Reserva{})
	gob.Register(models.User{})
	gob.Register(models.Restricao{})
	gob.Register(models.Livro{})
	gob.Register(map[string]bool{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData) // cria o channel para o envio de emails
	appConfig.MailChan = mailChan
	//mudar para true quando estiver em producao
	appConfig.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	mySession = scs.New()
	mySession.Lifetime = 24 * time.Hour
	mySession.Cookie.Persist = true
	mySession.Cookie.SameSite = http.SameSiteLaxMode
	mySession.Cookie.Secure = appConfig.InProduction
	appConfig.Session = mySession

	// conecta ao database
	log.Println("Conectando ao database...")
	db, err := myDriver.ConnectSQL(myDBLogin)
	if err != nil {
		log.Fatal("Não foi possível se conectar ao banco de dados.")
	}
	log.Println("Conectado ao database.")
	render.SetConfig(&appConfig)

	// tc é um mapa que armazena todos os templates html;
	// erro armazena possiveis erros que possam ocorrer no processamento dos templates html
	tc, erro := render.CreateTemplateCache()
	if erro != nil {
		log.Fatal("Nao foi possivel carregar os templates\n", erro)
		return nil, erro
	}

	// depois de carregados os templates, eles sao armazenados na variavel appConfig
	appConfig.TemplateCache = tc
	appConfig.UseCache = appConfig.InProduction // definido como false pois esta em desenvolvimento

	repo := handlers.NewRepo(&appConfig, db) // cria variavel do tipo handlers.Repository para as configs e o db poder ser utilizado neste package
	handlers.SetRepo(repo)                   // passa as configs para o pkg handlers
	render.SetConfig(&appConfig)             // passa as configs para o pkg render
	helpers.NewHelpers(&appConfig)           // passa as configs para o pkg helpers
	return db, nil
}
