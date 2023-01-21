package report

import (
	"testing"

	"github.com/cli/go-gh/pkg/auth"
)

func TestConfig_AuthToken(t *testing.T) {
	c := &Config{}
	token, _ := c.AuthToken("foo.com")
	token2, _ := auth.TokenForHost("foo.com")

	if token != token2 {
		t.Errorf("expected %s to equal %s", token, token2)
	}
}
