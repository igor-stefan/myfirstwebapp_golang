package dbrepo

import (
	"errors"
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

func (m *testPostgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation acrecenta uma linha na tabela 'Reserva' do db
func (m *testPostgresDBRepo) InsertReserva(res models.Reserva) (int, error) {
	return 1, nil
}

// InsertLivroRestricao insere uma nova restricao para determinado livro no db
func (m *testPostgresDBRepo) InsertLivroRestricao(r models.LivroRestricao) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID retorna true se existe disponibilidade
func (m *testPostgresDBRepo) SearchAvailabilityByDatesByRoomID(inicio, fim time.Time, livroID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms retorna um slice de livros que estao disponiveis para as datas especificados
func (m *testPostgresDBRepo) SearchAvailabilityForAllRooms(inicio, final time.Time) ([]models.Livro, error) {
	var livros []models.Livro
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
