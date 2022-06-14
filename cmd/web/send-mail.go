package main

import (
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() { // funcao anonima, sem nome e sem parametros
		for { // loop infinito
			msg := <-appConfig.MailChan
			enviaMsg(msg)
		}
	}() // para ser executada
}

func enviaMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// username
	// password
	// encyption

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)
	err = email.Send(client)
	if err != nil {
		errorLog.Println("err")
	} else {
		infoLog.Println("email enviado!")
	}
}
