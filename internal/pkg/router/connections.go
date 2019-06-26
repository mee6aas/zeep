package router

import (
	v1 "github.com/mee6aas/zeep/pkg/api/router/grpc/v1"
)

var (
	connections = make(map[string](chan v1.Message))
)
