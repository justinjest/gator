package json_parser

import (
	"testing"

	json_parser "github.com/justinjest/gator/internal/config"
)

func TestRead(t *testing.T) {
	config, err := json_parser.Read()
	if err != nil {
		t.Fatalf(`Errror unpacking config file = %q, %v`, config, err)
	}
}
