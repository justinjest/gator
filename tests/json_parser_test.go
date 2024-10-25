package json_parser

import (
	"testing"
)

func TestgetConfigPath(t *testing.T) {
	want := "~/gatorconfig.json"
	path, err := getConfigFilePath()
	if !want.MatchString(path) || err != nil {
		t.Fatalf(`~/gatorconfig.json = %q, %v, want match for %#q, nil`, path, err, want)
	}
}
