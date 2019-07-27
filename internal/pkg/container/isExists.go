package container

import (
	"context"
)

// IsExists checks if this container exists.
func (c Container) IsExists(ctx context.Context) (ok bool) {
	_, err := c.Inspect(ctx)
	ok = err == nil
	return
}
