package forms

import (
	"fmt"
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
			f.Errors.Add(campo, "Este campo deve ser preenchido.")
		}
	}
}

// Has checha se determinado form possui o campo nomeado 'field_name'
func (f *Form) Has(field_name string) bool {
	x := f.Get(field_name)
	if x == "" {
		f.Errors.Add(field_name, "Este campo não existe.")
		return false
	} else {
		return true
	}
}

// TamMin checa se a string possui tamanho mínimo especificado
func (f *Form) TamMin(nome_campo string, tam int) bool {
	x := f.Get(nome_campo)
	if len(x) < tam {
		f.Errors.Add(nome_campo, fmt.Sprintf("O nome deve conter no mínimo %d caracteres", tam))
		return false
	}
	return true
}

// IsEmail checha se os valores apresentados no campo especificado formam um email válido
func (f *Form) IsEmail(campo string) bool {
	if !govalidator.IsEmail(f.Get(campo)) {
		f.Errors.Add(campo, "Endereço de email inválido.")
		return false
	}
	return true
}
