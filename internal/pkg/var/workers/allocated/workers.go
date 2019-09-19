package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

var (
	//               usernmae   actName
	workers = make(map[string]map[string][]worker.Worker)
)
