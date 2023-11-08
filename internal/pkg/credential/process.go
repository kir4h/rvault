package credential

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultVersion = "1"
)

var (
	ErrNoCommandPassed = fmt.Errorf("expected at least a single command to execute")
)

type (
	Data struct {
		Version string
		Token   string
	}
)

func GetToken(command []string) (string, error) {
	if len(command) == 0 {
		return "", ErrNoCommandPassed
	}

	stdout := new(strings.Builder)
	stderr := new(strings.Builder)

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = []string{
		fmt.Sprintf("VAULT_ADDR=%s", viper.GetString("global.address")),
	}

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to exec process %v: %w, stdout: %s, stderr: %s", command, err, stdout.String(), stderr.String())
	}

	raw := stdout.String()

	var data Data
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse process output %s: %w", raw, err)
	}

	if data.Version != DefaultVersion {
		return "", fmt.Errorf("unexpected credentail process format version %s, expected %s", data.Version, DefaultVersion)
	}

	return data.Token, nil
}
