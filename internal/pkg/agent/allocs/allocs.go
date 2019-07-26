package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

var (
	//                taskID
	allocs = make(map[string][]worker.Worker)
)
