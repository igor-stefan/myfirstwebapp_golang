package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Form cria uma struct customizavel para os formularios e contém um objeto do tipo url.Values
type Form struct {
	url.Values
	Errors errors
}

// New inicializa uma struct do tipo Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid retorna true se não há erros, caso contrário retorna false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checa se o(s) campo(s) especificado(s) foi(ram) preenchido(s)
func (f *Form) Required(campos ...string) {
	for _, campo := range campos {
		value := f.Get(campo)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(campo, "Esse campo deve ser preenchido.")
		}
	}
}

// Has checha se determinado form possui um campo nomeado 'field_name'
func (f *Form) Has(field_name string, r *http.Request) bool {
	x := r.Form.Get(field_name)
	if x == "" {
		f.Errors.Add(field_name, "Este campo deve ser preenchido.")
		return false
	} else {
		return true
	}
}

// TamMin checa se a string possui tamanho mínimo especificado
func (f *Form) TamMin(nome_campo string, tam int, r *http.Request) bool {
	x := r.Form.Get(nome_campo)
	if len(x) < tam {
		f.Errors.Add(nome_campo, fmt.Sprintf("O nome deve conter no mínimo %d caracteres", tam))
		return false
	}
	return true
}

func (f *Form) IsEmail(campo string) {
	if !govalidator.IsEmail(f.Get(campo)) {
		f.Errors.Add(campo, "Endereço de email inválido.")
	}
}
