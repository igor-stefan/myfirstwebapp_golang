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
func (m *postgresDBRepo) InsertReserva(res models.Reserva) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()

	var returnedID int
	stmt := `insert into 
			reservas (nome, sobrenome, email, phone, data_inicio, data_final, livro_id, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.Nome,
		res.Sobrenome,
		res.Email,
		res.Phone,
		res.DataInicio,
		res.DataFinal,
		res.LivroID,
		time.Now(),
		time.Now(),
	).Scan(&returnedID)
	if err != nil {
		return -1, err
	}
	return returnedID, nil
}

// InsertLivroRestricao insere uma nova restricao para determinado livro no db
func (m *postgresDBRepo) InsertLivroRestricao(r models.LivroRestricao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()

	stmt := `insert into 
	livros_restricoes (data_inicio, data_final, id_livro, id_reserva, created_at, updated_at, id_restricao)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.DataInicio,
		r.DataFinal,
		r.LivroID,
		r.ReservaID,
		time.Now(),
		time.Now(),
		r.RestricaoID,
	)
	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID retorna true se existe disponibilidade
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(inicio, fim time.Time, livroID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()

	query := `select count(id)
		from livros_restricoes
		where id_livro = $1 and $2 <= data_final and $3 >= data_inicio;`

	var numRows int
	row := m.DB.QueryRowContext(ctx, query, livroID, inicio, fim)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// SearchAvailabilityForAllRooms retorna um slice de livros que estao disponiveis para as datas especificados
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(inicio, final time.Time) ([]models.Livro, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()

	var livros []models.Livro

	query := `select l.id_livro, l.nome_livro
	from livros l
	where l.id_livro not in
	(select id_livro from livros_restricoes lr where $1 <= lr.data_final and $2 >= lr.data_inicio);`

	rows, err := m.DB.QueryContext(ctx, query, inicio, final)
	if err != nil {
		return livros, err
	}

	for rows.Next() {
		var livro models.Livro
		err := rows.Scan(
			&livro.ID,
			&livro.NomeLivro,
		)
		if err != nil {
			return livros, err
		}
		livros = append(livros, livro)
	}

	if err = rows.Err(); err != nil {
		return livros, err
	}
	return livros, nil
}

// GetLivroByID busca no database o livro que possui o id especificado
func (m *postgresDBRepo) GetLivroByID(ID int) (models.Livro, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //cria um contexto para evitar que a solicitacao demore mais que 3 seg
	defer cancel()

	var livro models.Livro

	query := `select id_livro, nome_livro, created_at, updated_at from livros where id_livro = $1`
	row := m.DB.QueryRowContext(ctx, query, ID)
	err := row.Scan(
		&livro.ID,
		&livro.NomeLivro,
		&livro.CreatedAt,
		&livro.UpdatedAt,
	)
	if err != nil {
		return livro, err
	}
	return livro, nil
}
