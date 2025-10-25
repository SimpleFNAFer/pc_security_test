//go:build darwin
// +build darwin

package tester

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func CheckFWExists() (bool, error) {
	cmd := exec.Command("pfctl", "-s", "info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "error checking mac firewall")
	}
	return strings.Contains(string(output), "Status: Enabled"), nil
}

func CheckAVExists() (bool, error) {
	knownAV := []string{
		"/Applications/NortonSecurity.app",
		"/Applications/AvastSecurity.app",
		"/Applications/Malwarebytes.app",
	}
	for _, path := range knownAV {
		if err := exec.Command("test", "-d", path).Run(); err == nil {
			return true, nil
		}
	}
	return false, nil
}
