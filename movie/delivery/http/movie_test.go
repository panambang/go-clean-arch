package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/bxcodec/go-clean-arch/domain/mocks"
)

func TestFetch(t *testing.T) {
	var mockMovie domain.Movie
	err := faker.FakeData(&mockMovie)
	assert.NoError(t, err)
	mockUCase := new(mocks.MovieUsecase)
	mockListMovie := make([]domain.Movie, 0)
	mockListMovie = append(mockListMovie, mockMovie)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListMovie, "10", nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/movie?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := movieHttp.MovieHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchMovie(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.MovieUsecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/movie?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := movieHttp.MovieHandler{
		AUsecase: mockUCase,
	}
	err = handler.FetchMovie(c)
	require.NoError(t, err)

	responseCursor := rec.Header().Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockMovie domain.Movie
	err := faker.FakeData(&mockMovie)
	assert.NoError(t, err)

	mockUCase := new(mocks.MovieUsecase)

	num := int(mockMovie.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockMovie, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/movie/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("movie/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := movieHttp.MovieHandler{
		AUsecase: mockUCase,
	}
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockMovie := domain.Movie{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockMovie := mockMovie
	tempMockMovie.ID = 0
	mockUCase := new(mocks.MovieUsecase)

	j, err := json.Marshal(tempMockMovie)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.Movie")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/movie", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/movie")

	handler := movieHttp.MovieHandler{
		AUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockMovie domain.Movie
	err := faker.FakeData(&mockMovie)
	assert.NoError(t, err)

	mockUCase := new(mocks.MovieUsecase)

	num := int(mockMovie.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/movie/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("movie/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := movieHttp.MovieHandler{
		AUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUCase.AssertExpectations(t)

}
