package repository

import (
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

// DataBaseRepo é uma interface que permite que as funcionalidades armazenadas sejam utilizadas em outros pkgs (handlers, por exemplo)
type DataBaseRepo interface {
	AllUsers() bool
	InsertReserva(res models.Reserva) (int, error)
	InsertLivroRestricao(r models.LivroRestricao) error
	SearchAvailabilityByDatesByLivroID(inicio, fim time.Time, livroID int) (bool, error)
	SearchAvailabilityForAllLivros(inicio, final time.Time) ([]models.Livro, error)
	GetLivroByID(ID int) (models.Livro, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Autenticar(email, senhaFornecida string) (int, string, error)
}
