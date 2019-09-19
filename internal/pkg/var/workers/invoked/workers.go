package workers

import (
	"github.com/mee6aas/zeep/internal/pkg/worker"
)

var (
	//		   		  Cont.IP
	workers = make(map[string]*worker.Worker)
)
