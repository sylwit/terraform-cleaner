package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindTfModules(t *testing.T) {
	t.Run("must found 3 modules", func(t *testing.T) {
		path := "./testdata"

		result, err := findTfModules(path)
		assert.Equal(t, err, nil)
		assert.Equal(t, len(result), 3)

		assert.Contains(t, result, "testdata")
		assert.Contains(t, result, "testdata/tf")
		assert.Contains(t, result, "testdata/tf/1")
	})
}

func TestFindVariablesUsage(t *testing.T) {
	t.Run("should find all variables", func(t *testing.T) {
		path := "./testdata/tf"

		out, err := findVariablesUsage(path)
		assert.Equal(t, err, nil)
		assert.Equal(t, 1, out.variables["name"])
		assert.Equal(t, 1, out.variables["region"])
		assert.Equal(t, 1, out.variables["instance_ids"])
		assert.Equal(t, 0, out.variables["legacy"])
	})
}
