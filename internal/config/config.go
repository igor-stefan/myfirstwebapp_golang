package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig é uma struct que armazena as configuracoes da aplicacao
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}
