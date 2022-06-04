package helpers

import (
	"fmt"
	"strings"
	"time"
)

// ConvStr2Time converte uma string em formato de data para uma variavel do tipo time.Time
func ConvStr2Time(layout, str2Bconv string) (time.Time, error) {
	str2Bconv = strings.ReplaceAll(str2Bconv, "/", "-")
	layout = strings.ReplaceAll(layout, "/", "-")
	var strConverted time.Time
	var myerr error = nil
	strConverted, myerr = time.Parse(layout, str2Bconv)
	if myerr != nil {
		myerr = fmt.Errorf("nao foi possivel fazer a conversao de string para formato de tempo\n%s", myerr)
		return strConverted, myerr
	}
	return strConverted, nil
}
