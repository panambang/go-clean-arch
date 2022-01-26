package movie_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/movie/repository"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockMovies := []domain.Movie{
		domain.Movie{
			ID: 1, Title: "title 1", Content: "content 1",
			Logmovie: domain.Logmovie{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.Movie{
			ID: 2, Title: "title 2", Content: "content 2",
			Logmovie: domain.Logmovie{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "logmovie_id", "updated_at", "created_at"}).
		AddRow(mockMovies[0].ID, mockMovies[0].Title, mockMovies[0].Content,
			mockMovies[0].Logmovie.ID, mockMovies[0].UpdatedAt, mockMovies[0].CreatedAt).
		AddRow(mockMovies[1].ID, mockMovies[1].Title, mockMovies[1].Content,
			mockMovies[1].Logmovie.ID, mockMovies[1].UpdatedAt, mockMovies[1].CreatedAt)

	query := "SELECT id,title,content, logmovie_id, updated_at, created_at FROM movie WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := movieMysqlRepo.NewMysqlMovieRepository(db)
	cursor := repository.EncodeCursor(mockMovies[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "logmovie_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, logmovie_id, updated_at, created_at FROM movie WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := movieMysqlRepo.NewMysqlMovieRepository(db)

	num := int64(5)
	anMovie, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anMovie)
}

func TestStore(t *testing.T) {
	now := time.Now()
	ar := &domain.Movie{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Logmovie: domain.Logmovie{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  movie SET title=\\? , content=\\? , logmovie_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Logmovie.ID, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := movieMysqlRepo.NewMysqlMovieRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "logmovie_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id,title,content, logmovie_id, updated_at, created_at FROM movie WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := movieMysqlRepo.NewMysqlMovieRepository(db)

	title := "title 1"
	anMovie, err := a.GetByTitle(context.TODO(), title)
	assert.NoError(t, err)
	assert.NotNil(t, anMovie)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM movie WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := movieMysqlRepo.NewMysqlMovieRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	ar := &domain.Movie{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Logmovie: domain.Logmovie{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE movie set title=\\?, content=\\?, logmovie_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Logmovie.ID, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := movieMysqlRepo.NewMysqlMovieRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
