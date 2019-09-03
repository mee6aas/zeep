package v1

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mee6aas/zeep/internal/pkg/agent/allocs"
)

// ResolveNameFromIP is invoked when requests needs to handle unamed request.
func (h Handle) ResolveNameFromIP(
	_ context.Context,
	ip string,
) (n string, e error) {
	var (
		ok bool
	)

	n, ok = allocs.ResolveNameFromIP(ip)
	if !ok {
		e = errors.New("Not found")
	}

	return
}
