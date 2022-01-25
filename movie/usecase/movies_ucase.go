package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type movieUsecase struct {
	movieRepo      domain.MovieRepository
	contextTimeout time.Duration
}

// NewMovieUsecase will create new an movieUsecase object representation of domain.MovieUsecase interface
func NewMovieUsecase(a domain.MovieRepository, timeout time.Duration) domain.MovieUsecase {
	return &movieUsecase{
		movieRepo:      a,
		contextTimeout: timeout,
	}
}

func (a *movieUsecase) Fetch(c context.Context, cursor string, searchword string) (res []domain.Movies, nextCursor string, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.movieRepo.Fetch(ctx, cursor, searchword)
	if err != nil {
		return nil, "", err
	}

	return res, nextCursor, nil
}

func (a *movieUsecase) GetByID(c context.Context, id string) (res domain.Movies, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.movieRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return res, nil
}
