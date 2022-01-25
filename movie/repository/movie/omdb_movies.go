package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bxcodec/go-clean-arch/domain"
)

type omdbAPIRepository struct {
	APIKey string
}

const (
	omdbBaseURL = "http://www.omdbapi.com/"
)

// NewMysqlMovieRepository will create an object that represent the movie.Repository interface
func NewMysqlMovieRepository(APIKey string) domain.MovieRepository {
	return &omdbAPIRepository{APIKey}
}

func (m *omdbAPIRepository) Fetch(ctx context.Context, page string, searchword string) (res []domain.Movies, nextCursor string, err error) {
	var client = &http.Client{}
	var movies domain.SearchResult

	url := fmt.Sprintf("%s?apikey=%s&s=%s&page=%s", omdbBaseURL, m.APIKey, searchword, page)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&movies)
	if err != nil {
		return
	}

	return movies.Search, "1", err
}

func (m *omdbAPIRepository) GetByID(ctx context.Context, imdbID string) (res domain.Movies, err error) {
	var client = &http.Client{}
	var movies domain.Movies

	url := fmt.Sprintf("%s?apikey=%s&i=%s&plot=full", omdbBaseURL, m.APIKey, imdbID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&movies)
	if err != nil {
		return
	}

	return movies, err
}
