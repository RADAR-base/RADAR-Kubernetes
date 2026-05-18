package output_test

import (
	"testing"

	"github.com/RADAR-base/RADAR-Kubernetes/cli/pkg/output"
)

func TestModeConstants(t *testing.T) {
	if output.Human != "human" {
		t.Fatal("expected Human == 'human'")
	}
	if output.JSON != "json" {
		t.Fatal("expected JSON == 'json'")
	}
}
