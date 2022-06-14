package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

// AppConfig é uma struct que armazena as configuracoes da aplicacao
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
