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
		res *v1.InvokeResponse
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if res, err = testClient.Invoke(ctx, &v1.InvokeRequest{
		Username: "Rick",
		ActName:  "Zeep",
		Arg:      "Peace",
	}); err != nil {
		t.Fatalf("failed to request to invoke: %v", err)
	}

	if res.GetResult() != "Zeep by Rick with Peace" {
		t.Fatalf("Undesired result")
	}
}
