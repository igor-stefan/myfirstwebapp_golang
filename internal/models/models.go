package models

import "time"

// User é o modelo para a table User no db
type User struct {
	Id          int
	Nome        string
	SobreNome   string
	Email       string
	Senha       string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Livro é o modelo para a table Livro no db
type Livro struct {
	ID          int
	NomeLivro   string
	Emprestimos int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Restricao é o modelo para a table Restricao no db
type Restricao struct {
	ID            int
	NomeRestricao string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Reservas é o modelo para a table Reservas no db
type Reserva struct {
	ID         int
	Nome       string
	Sobrenome  string
	Email      string
	Phone      string
	Obs        string
	DataInicio time.Time
	DataFinal  time.Time
	Processada int
	LivroID    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Livro      Livro
}

// LivrosRestricao é o modelo para a table LivrosRestricao no db
type LivroRestricao struct {
	ID          int
	DataInicio  time.Time
	DataFinal   time.Time
	LivroID     int
	ReservaID   int
	RestricaoID int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Livro       Livro
	Reserva     Reserva
	Restricao   Restricao
}

// MailData possui as informacoes para envio de email
type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
