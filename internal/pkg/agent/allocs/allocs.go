package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

var (
	//                actID
	allocs = make(map[string][]worker.Worker)
)
