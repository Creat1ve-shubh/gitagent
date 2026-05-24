package tools

import "context"

type Metadata struct {
	ConcurrencySafe bool
	ReadOnly        bool
	Destructive     bool
	MaxResultChars  int
}

type Tool interface {
	Name() string
	Description() string
	Schema() map[string]any
	Metadata() Metadata
	Execute(ctx context.Context, args map[string]any) (string, error)
}

func DefaultMetadata() Metadata {
	return Metadata{ConcurrencySafe: false, ReadOnly: false, Destructive: false, MaxResultChars: 50000}
}
