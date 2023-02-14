package credential

import (
	"errors"
	"strings"
	"testing"
)

func TestGetToken_SingleCommand(t *testing.T) {
	_, err := GetToken([]string{"test"})
	if !strings.Contains(err.Error(), "exit status 1") {
		t.Errorf("GetToken got error = %v, want exit status 1", err)
	}
}

func TestGetToken_NoCommand(t *testing.T) {
	_, err := GetToken([]string{})
	if !errors.Is(err, NoCommandPassed) {
		t.Errorf("GetToken got error = %v, want NoCommandPassed", err)
	}
}

func TestGetToken_HappyPath(t *testing.T) {
	token, err := GetToken([]string{"sh", "-c", `echo '{"Version":"1","Token":"123"}'`})
	if err != nil {
		t.Errorf("GetToken got error = %v, want nil", err)
	}
	if token != "123" {
		t.Errorf("GetToken got token = %v, want 123", token)
	}
}
