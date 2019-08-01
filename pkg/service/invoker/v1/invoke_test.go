package v1_test

import (
	"context"
	"testing"
	"time"

	v1 "github.com/mee6aas/zeep/pkg/api/invoker/v1"
)

func TestInvoke(t *testing.T) {
	var (
		err error
		ctx context.Context
		cc  context.CancelFunc
		res *v1.InvokeResponse
	)

	ctx, cc = context.WithTimeout(context.Background(), time.Second)
	defer cc()

	if res, err = testClient.Invoke(ctx, &v1.InvokeRequest{
		Target: &v1.Activity{
			Name:  "Zeep",
			Label: "C-137",
		},
	}); err != nil {
		t.Fatalf("failed to request to invoke: %v", err)
	}

	if res.GetResult() != "Zeep@C-137" {
		t.Fatalf("Undesired result")
	}
}
