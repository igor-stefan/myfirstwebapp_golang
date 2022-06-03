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
	SearchAvailabilityByDatesByRoomID(inicio, fim time.Time, livroID int) (bool, error)
	SearchAvailabilityForAllRooms(inicio, final time.Time) ([]models.Livro, error)
}
