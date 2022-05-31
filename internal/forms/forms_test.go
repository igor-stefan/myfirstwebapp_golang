package forms

import (
	"fmt"
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
	form = New(r.PostForm)     // atualiza o form com os novos valores

	// checa se estes campos, que são obrigatórios, foram preenchidos
	// todos os campos passados para a funcao required sao obrigatorios
	// se nao forem preeenchidos, um erro é adicionado, entao basta checar se o formulario é válido
	form.Required("a", "b", "c")

	// checa se há erros adicionados
	// esse teste deve falhar, pois nao devem ter sido adicionado erros
	if !form.Valid() {
		t.Error("indica que não possui campos obrigatórios preenchidos, entretanto, possui")
	}
}

func TestForm_Has(t *testing.T) {
	dadosPostados := url.Values{}      // cria uma variavel do tipo url.values
	form := New(dadosPostados)         // cria um form com os dados passados pela requisicao (vazio)
	to_test := []string{"a", "b", "c"} // nomes dos campos que devem ser testados
	for _, campo := range to_test {    // testa cada um dos campos (vazio)
		has := form.Has(campo)
		if has {
			t.Error("o teste mostra que existe o campo no formulario, entretanto, NAO existe")
		}
	}

	for _, cmp := range to_test {
		dadosPostados.Add(cmp, "valor_"+cmp) // adiciona dados ao campo "cmp"
	}
	form = New(dadosPostados) // atualiza o form com os novos valores (pense como renderização do form)
	for _, campo := range to_test {
		has := form.Has(campo)
		if !has {
			t.Error("o teste mostra que NAO existe o campo no formulario, entretanto, existe")
		}
	}
}

func TestForm_TamMin(t *testing.T) {
	dadosPostados := url.Values{}               // cria uma variavel do tipo url values vazia
	form := New(dadosPostados)                  // adiciona ao form os dados em url values (vazio, no caso)
	minLength := 3                              // define o menor tamanho para ser checado nos testes
	hasMinLength := form.TamMin("a", minLength) // realiza o teste a vazio
	if hasMinLength {
		t.Error("o teste mostra que o campo tem os dados com o tam. minimo, porém NAO tem")
	}
	to_test := []string{"a", "ab", "abc", "abcd", "abcde", "abcdef"} // cria os vlaores para serem testados
	for _, campo := range to_test {
		//adiciona valores a variavel do tipo url values para serem adicionados ao form
		dadosPostados.Add("chave_"+campo, campo)
	}
	form = New(dadosPostados) // atualiza os valores do form (antes era vazio)
	for j, campo := range to_test {
		hasMinLength = form.TamMin("chave_"+campo, minLength)
		if hasMinLength && j < 2 {
			t.Error("o teste mostra que o campo tem os dados com o tam. min., entretanto, NAO tem")
		} else if !hasMinLength && j >= 2 {
			t.Error("o teste mostra que o campo NAO tem os dados com o tam. min., entretanto, tem")
		}
	}
}

func TestForm_IsEmail(t *testing.T) {
	dadosPostados := url.Values{} // cria uma variavel do tipo url values vazia
	to_test := []string{"a@a", "b@b.", "a#b.com", "asdasjd@@x.com", "d@d.com", "c@c.c", "ppppppp@zzzzz.com"}
	for j, val := range to_test {
		dadosPostados.Add(fmt.Sprint(j), val)
	}
	form := New(dadosPostados) // adiciona ao form os dados em url values (vazio, no caso)
	for j := range to_test {
		isEmail := form.IsEmail(fmt.Sprint(j)) // chama a funcao a ser testada
		if isEmail && j < 4 {
			t.Error("o teste mostra que o valor é um email, entretanto, NAO é. teste", j+1)
		} else if !isEmail && j > 3 {
			t.Error("o teste mostra que o valor NAO é um email, entretanto, é. teste", j+1)
		}
	}
}
