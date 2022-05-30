package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/qualquer", nil)
	form := New(r.PostForm) // cria um form com os dados passados para a requisicao

	isValid := form.Valid() //testa a funcao
	if !isValid {
		t.Error("o formulario foi dado como invalido, mas deveria ser validado")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/qualquer", nil) // cria uma requisicao para o teste ser executado
	form := New(r.PostForm)                            // cria um form com os dados passados para a requisicao (no caso, nenhum dado)

	form.Required("a", "b", "c") // chama a funcao required do pkg forms com os campos a,b,c (diz que estes campos sao obrigatorios e checa seu preenchimento)
	// testa se o formulario é valido (se nao houve erros adicionados)
	if form.Valid() { // esse teste deve falhar, caso sucesso, então ocorreu um erro, pois o form é inválido
		t.Error("o formulario foi dado como valido, mas os campos obrigatórios deveriam estar preenchidos")
	}

	dadosPostados := url.Values{} // cria uma variavel do tipo url.values
	dadosPostados.Add("a", "a")   // adiciona dados ao campo "a"
	dadosPostados.Add("b", "b")   // adiciona dados ao campo "b"
	dadosPostados.Add("c", "c")   // adiciona dados ao campo "c"

	r = httptest.NewRequest("POST", "/show", nil) //cria uma nova request

	r.PostForm = dadosPostados //agora os dados em postForm da nova request recebem os url.values definidos anteriormente
	form = New(r.PostForm)     // cria um novo form

	// checa se estes campos, que são obrigatórios, foram preenchidos
	// todos os campos passados para a funcao required sao obrigatorios
	// se nao forem preeenchidos, um erro é adicionado
	form.Required("a", "b", "c")

	// checa se há erros adicionados
	// esse teste deve falhar, pois nao devem ter sido adicionado erros
	if !form.Valid() {
		t.Error("indica que não possui campos obrigatórios preenchidos, quando, na verdade, possui")
	}
}
