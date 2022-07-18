package dbrepo

import (
	"errors"
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/helpers"
	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

func (m *testPostgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation acrecenta uma linha na tabela 'Reserva' do db
func (m *testPostgresDBRepo) InsertReserva(res models.Reserva) (int, error) {
	if res.LivroID == -1 {
		return res.LivroID, errors.New("numero do livro invalido, nao foi possivel inserir a reserva no db")
	}
	return res.LivroID, nil
}

// InsertLivroRestricao insere uma nova restricao para determinado livro no db
func (m *testPostgresDBRepo) InsertLivroRestricao(r models.LivroRestricao) error {
	if r.LivroID == -2 {
		return errors.New("numero do livro invalido, nao foi possivel inserir a restricao do livro no db")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID retorna true se existe disponibilidade
func (m *testPostgresDBRepo) SearchAvailabilityByDatesByLivroID(inicio, fim time.Time, livroID int) (bool, error) {
	if livroID == -1 {
		return false, nil
	} else if livroID == -2 {
		return false, errors.New("erro ao conectar-se ao db")
	}
	return true, nil
}

// SearchAvailabilityForAllLivros retorna um slice de livros que estao disponiveis para as datas especificados
func (m *testPostgresDBRepo) SearchAvailabilityForAllLivros(inicio, final time.Time) ([]models.Livro, error) {
	var layout string = "02-01-2006"
	var livros []models.Livro
	if time, _ := helpers.ConvStr2Time(layout, "01-01-3000"); inicio == time { // nao retorna livros, ano 3000
		return livros, nil
	} else if time, _ := helpers.ConvStr2Time(layout, "01-01-1999"); inicio == time { // retorna erro, ano 1999
		return livros, errors.New("nao foi possivel processar os dados com as informacoes recebidas")
	}
	livros = append(livros, models.Livro{
		ID:        1000,
		NomeLivro: "Ratos de cemitério",
	})
	return livros, nil
}

// GetLivroByID busca no database o livro que possui o id especificado
func (m *testPostgresDBRepo) GetLivroByID(ID int) (models.Livro, error) {
	var livro models.Livro
	if ID > 7 {
		return livro, errors.New("nao foi encontrado livro com ID especificado")
	}
	return livro, nil
}

// GetUserById retorna dados de um usuário especificado pelo seu ID
func (m *testPostgresDBRepo) GetUserById(id int) (models.User, error) {
	return models.User{}, nil
}

func (m *testPostgresDBRepo) UpdateUser(u models.User) error {
	return nil
}

// Autenticar retorna id e senha em caso de sucesso na autenticacao do usuário
func (m *testPostgresDBRepo) Autenticar(email, senhaFornecida string) (int, string, error) {
	return -10, "umasenhaqualquer", nil
}

// AllReservas retorna todas as reservas presentes no db
func (m *testPostgresDBRepo) AllReservas() ([]models.Reserva, error) {
	var reservas []models.Reserva
	return reservas, nil
}
