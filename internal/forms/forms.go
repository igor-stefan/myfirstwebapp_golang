package forms

import (
	"net/http"
	"net/url"
)

//Form cria uma struct customizavel para os formularios e contém um objeto do tipo url.Values
type Form struct {
	url.Values
	Errors errors
}

//New inicializa uma struct do tipo Form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checha se determinado form possui um campo nomeado 'field_name'
func (f *Form) Has(field_name string, r *http.Request) bool {
	x := r.Form.Get(field_name)
	if x == "" {
		return false
	} else {
		return true
	}
}
