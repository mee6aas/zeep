package assigns_test

import (
	"testing"
	"time"

	"github.com/mee6aas/zeep/internal/pkg/agent/assigns"
)

func TestAssigns(t *testing.T) {
	id, c := assigns.Add("0.0.0.0", "Jerry")

	t.Logf("id %v", id)

	go func() {
		time.Sleep(time.Millisecond * 10)

		assigns.Report(id, &struct{}{})
	}()

	select {
	case <-c:
	case <-time.After(time.Millisecond * 20):
		t.Fatal("Expected to get result")
	}
}
