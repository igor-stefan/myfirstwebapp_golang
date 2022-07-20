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

function callbackJsonAfterPickDate(dadosForm, id) { // checa se o livro que possui o 'id' especificado está disponivel e alerta o usuário
    //const dadosForm = new FormData(form) // adquire os dados do formulario de uma maneira que easily construct a set of key/value pairs representing form fields and their values
   
}