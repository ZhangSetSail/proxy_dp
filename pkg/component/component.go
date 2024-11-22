package component

import "context"

type Component interface {
	Start(ctx context.Context) error
	CloseHandle()
}
