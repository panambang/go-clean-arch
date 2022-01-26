package mysql

import (
	"context"
	"database/sql"

	"github.com/bxcodec/go-clean-arch/domain"
)

type mysqlLogmovieRepo struct {
	DB *sql.DB
}

// NewMysqlLogmovieRepository will create an implementation of logmovie.Repository
func NewMysqlLogmovieRepository(db *sql.DB) domain.LogmovieRepository {
	return &mysqlLogmovieRepo{
		DB: db,
	}
}

func (mm *mysqlLogmovieRepo) Store(ctx context.Context, m *domain.Movies) (err error) {
	query := `INSERT movies SET title=? , imdbID=? , year=?, released=? , imdbRating=?`
	stmt, err := mm.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, m.Title, m.ID, m.Year, m.Released, m.ImdbRating)
	if err != nil {
		return
	}
	_, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}
