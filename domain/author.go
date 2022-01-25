package domain

import "context"

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	Store(ctx context.Context, m *Movies) error
}
