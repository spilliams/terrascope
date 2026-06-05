package hcl

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDependencyGraph(t *testing.T) {
	type test struct {
		moduleDir     string
		expectedGraph string
	}

	tests := []test{
		{
			moduleDir: "../../fixtures/roots/listed-resource",
			expectedGraph: `digraph G {
	"var.qty"->"random_string.this";
	"random_string.this";
	"var.qty";

}`,
		},
		{
			moduleDir: "../../fixtures/roots/mapped-resource",
			expectedGraph: `digraph G {
	"var.keys"->"random_string.this";
	"random_string.this";
	"var.keys";

}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.moduleDir, func(t *testing.T) {
			parser := NewModule(logrus.StandardLogger())
			if err := parser.ParseModuleDirectory(tc.moduleDir); err != nil {
				t.Error(err)
			}

			actualGraph, err := parser.DependencyGraph()
			if err != nil {
				t.Error(err)
			}

			if strings.TrimSpace(actualGraph) != strings.TrimSpace(tc.expectedGraph) {
				t.Logf("Expected: %s", tc.expectedGraph)
				t.Logf("Actual:   %s", actualGraph)
				t.Error("Actual graph did not match expected graph.")
			}
		})
	}
}
