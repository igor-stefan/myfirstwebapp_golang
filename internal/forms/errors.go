package forms

type errors map[string][]string

//Add adiciona uma mensagem de erro para determinado campo
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//Get retorna a primeira msg de erro
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
