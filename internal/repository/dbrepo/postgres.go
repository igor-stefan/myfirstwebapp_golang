package dbrepo

import (
	"context"
	"time"

	"github.com/igor-stefan/myfirstwebapp_golang/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation acrecenta uma linha na tabela 'Reserva' do db
func (m *postgresDBRepo) InsertReserva(res models.Reserva) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()
	stmt := `insert into 
			reservas (nome, sobrenome, email, phone, data_inicial, data_final, livro_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.Nome,
		res.Sobrenome,
		res.Email,
		res.Phone,
		res.DataInicio,
		res.DataFinal,
		res.LivroID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
