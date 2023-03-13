package terraform

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListTfModules(t *testing.T) {
	t.Run("must found 3 modules", func(t *testing.T) {
		path := "./testdata"

		result, err := ListTfModules(path)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(result), 3)

		assert.Contains(t, result, "testdata")
		assert.Contains(t, result, "testdata/tf")
		assert.Contains(t, result, "testdata/tf/1")
	})
}

func TestNewModuleUsage(t *testing.T) {
	t.Run("should init and process ModuleUsage", func(t *testing.T) {
		path := "./testdata/tf"

		moduleUsage, err := NewModuleUsage(path)
		assert.Equal(t, err, nil)
		assert.Equal(t, 4, len(moduleUsage.Variables))
		assert.Equal(t, 3, len(moduleUsage.Locals))

		assert.Equal(t, 1, moduleUsage.Variables["name"])
		assert.Equal(t, 1, moduleUsage.Variables["region"])
		assert.Equal(t, 1, moduleUsage.Variables["instance_ids"])
		assert.Equal(t, 0, moduleUsage.Variables["legacy"])

		assert.Equal(t, 1, moduleUsage.Locals["tags"])
		assert.Equal(t, 0, moduleUsage.Locals["dummy"])
		assert.Equal(t, 0, moduleUsage.Locals["dummy2"])
	})
}
