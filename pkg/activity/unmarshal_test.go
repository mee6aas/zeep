package activity_test

import (
	"fmt"
	"testing"

	"gotest.tools/assert"

	"github.com/mee6aas/zeep/pkg/activity"
)

// TODO: more testcases

func TestUnmarshalFromFile(t *testing.T) {
	{
		const (
			testActivity    = "./testdata/withoutDeps.json"
			expectedRuntime = "mee6aas/runtime-nodejs"
		)

		a, e := activity.UnmarshalFromFile(testActivity)
		assert.NilError(t, e,
			fmt.Sprintf("Expected to unmarshal %s\n", testActivity))
		assert.Equal(t, a.Runtime, expectedRuntime,
			fmt.Sprintf("Expected that the runtime is %s\n", expectedRuntime))
		assert.Equal(t, len(a.Dependencies), 0,
			"Expected that the runtime has no dependencies\n")
	}

	{
		const testActivity = "./testdata/withDeps.json"
		var expectedDpes = map[string]activity.Dep{
			"meeseeks1": activity.Dep{},
			"meeseeks2": activity.Dep{Outflow: "optional"},
		}

		a, e := activity.UnmarshalFromFile(testActivity)
		assert.NilError(t, e,
			fmt.Sprintf("Expected to unmarshal %s\n", testActivity))
		assert.DeepEqual(t, a.Dependencies, expectedDpes)
	}
}
