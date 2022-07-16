package models

import "github.com/igor-stefan/myfirstwebapp_golang/internal/forms"

// TemplateData guarda os dados enviados dos handlers para os templates html
type TemplateData struct {
	StringMap   map[string]string
	IntMap      map[string]int
	FloatMap    map[string]float32
	Data        map[string]interface{}
	CSRFToken   string
	Flash       string
	Warning     string
	Error       string
	Form        *forms.Form
	Autenticado int
}
