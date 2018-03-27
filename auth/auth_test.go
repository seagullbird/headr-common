package auth_test

import (
	"github.com/seagullbird/headr-common/auth"
	"testing"
)

func TestLogin(t *testing.T) {
	if auth.Login() == "" {
		t.Fatal()
	}
}
