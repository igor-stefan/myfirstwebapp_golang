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
