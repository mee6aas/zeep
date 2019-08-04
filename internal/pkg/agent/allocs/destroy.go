package allocs

import "context"

// Destroy removes allocated workers.
func Destroy(ctx context.Context) (e error) {
	for _, es := range allocs {
		for _, ws := range es {
			for _, w := range ws {
				e = w.Remove(ctx)
			}
		}
	}

	return
}
