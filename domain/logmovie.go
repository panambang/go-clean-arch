package domain

import "context"

// LogmovieRepository represent the logmovie's repository contract
type LogmovieRepository interface {
	Store(ctx context.Context, m *Movies) error
}
