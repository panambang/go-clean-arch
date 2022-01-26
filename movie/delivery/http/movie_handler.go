package http

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// MovieHandler  represent the httphandler for movie
type MovieHandler struct {
	MUsecase domain.MovieUsecase
	LogRepo  domain.LogmovieRepository
}

// NewMovieHandler will initialize the movies/ resources endpoint
func NewMovieHandler(e *echo.Echo, us domain.MovieUsecase, lr domain.LogmovieRepository) {
	handler := &MovieHandler{
		MUsecase: us,
		LogRepo:  lr,
	}
	e.GET("/movies", handler.FetchMovie)
	e.GET("/movies/:id", handler.GetByID)
}

// FetchMovie will fetch the movie based on given params
func (a *MovieHandler) FetchMovie(c echo.Context) error {
	searchword := c.QueryParam("searchword")
	page := c.QueryParam("paginatioon")
	ctx := c.Request().Context()

	listAr, _, err := a.MUsecase.Fetch(ctx, page, searchword)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, listAr)
}

// GetByID will get movie by given id
func (a *MovieHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()

	art, err := a.MUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.LogRepo.Store(ctx, &art)
		if err != nil {
			logrus.Warnf(err.Error())

		}
	}()

	wg.Wait()

	return c.JSON(http.StatusOK, art)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
