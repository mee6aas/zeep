package activity_test

import (
	"testing"

	"github.com/mee6aas/zeep/pkg/activity"
)

// TODO: more testcases

func TestUnmarshalFromFile(t *testing.T) {
	p := "./testdata/valid/activity.json"

	a, e := activity.UnmarshalFromFile(p)

	if e != nil {
		t.Fatalf("Expected to unmarshal %s: %v", p, e)
	}

	if n := a.Name; n != "meeseeks" {
		t.Fatalf("Expected that the name is meeseeks but %s", n)
	}

	if r := a.Runtime; r != "mee6aas/nodejs" {
		t.Fatalf("Expected that the runtime is mee6aas/nodejs but %s", r)
	}
}

func TestUnmarshalFromDir(t *testing.T) {
	p := "./testdata/valid"

	a, e := activity.UnmarshalFromDir(p)

	if e != nil {
		t.Fatalf("Expected to unmarshal %s: %v", p, e)
	}

	if n := a.Name; n != "meeseeks" {
		t.Fatalf("Expected that the name is meeseeks but %s", n)
	}

	if r := a.Runtime; r != "mee6aas/nodejs" {
		t.Fatalf("Expected that the runtime is mee6aas/nodejs but %s", r)
	}
}
