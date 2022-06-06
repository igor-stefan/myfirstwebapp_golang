function processarData(date) {
    ano = date.split('-')[0]
    mes = date.split('-')[1]
    dia = date.split('-')[2]
    ret = dia + '-' + mes + '-' + ano
    return ret
}

function processarString(str) {
    dia = str.split('-')[0]
    mes = str.split('-')[1]
    ano = str.split('-')[2]
    ret = new Date(ano + '-' + mes + '-' + dia + 'T05:00:00.000Z')
    return ret
}

function callbackJsonAfterPickDate(nomeForm, id) { // checa se o livro que possui o 'id' especificado está disponivel e alerta o usuário
    const form = document.getElementById(nomeForm); // formulario com as datas
    console.log("elem = ", form)
    const dadosForm = new FormData(form) // adquire os dados do formulario de uma maneira que easily construct a set of key/value pairs representing form fields and their values
    dadosForm.append("csrf_token", "{{ .CSRFToken }}"); // acrescenta o token
    dadosForm.append("id_livro", id); // acrescenta o id do livro
    fetch("/catalogo-json", { // realiza um fetch para a rota especifica
        method: "post", // o handler PostCatalogoJson lida com esse fetch
        body: dadosForm
    })
        .then(respostaFetch => respostaFetch.json()) // recebe a resposta e a converte para json
        .then(dados => { //dados -> variavel que contém o arquivo json
            if (dados.ok) {
                let msg = `<p>Livro Disponível!</p>
                <p><a 
                    href="/reservar-livro?id=${dados.livroID}&di=${dados.dataInicio}&df=${dados.dataFinal}"
                    class="btn btn-primary"> Reservar Agora!
                </a></p>
                `
                alertar.custom({
                    icon: "success",
                    title: "Ótima Notícia!",
                    showConfirmButton: false,
                    msg: msg,
                })
            } else {
                alertar.erro({
                    msg: "Livro indisponível!",
                    icon: "error",
                    text: "Tente outro período",
                    footer: "",
                })
            }
        })
}