package allocs

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

var (
	//               usernmae   actName
	allocs = make(map[string]map[string][]worker.Worker)
)
