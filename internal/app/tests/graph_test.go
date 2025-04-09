package tests

import (
	"regexp"
	"testing"
)

func TestReadFile(t *testing.T) {
	type testCase struct {
		filename string
		properties []Property

		// Expected Values
		expected string
		err error
	}

	t.Run("Valid file read values", func(t *testing.T) {
		tests := []testCase {}
	})
}
