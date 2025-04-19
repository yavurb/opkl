package opkl

import "context"

type Client interface {
	ReadSecret(ctx context.Context, ref string) (string, error)
}

type option func(*opklReader) error
