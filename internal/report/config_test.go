package report

import (
	"testing"

	"github.com/cli/go-gh/v2/pkg/auth"
)

func TestConfig_ActiveToken(t *testing.T) {
	c := &Config{}
	token, _ := c.ActiveToken("foo.com")
	token2, _ := auth.TokenForHost("foo.com")

	if token != token2 {
		t.Errorf("expected %s to equal %s", token, token2)
	}
}
