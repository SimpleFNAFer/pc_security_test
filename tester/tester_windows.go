//go:build windows
// +build windows

package tester

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/yusufpapurcu/wmi"
)

func CheckFWExists() (bool, error) {
	cmd := exec.Command("sc", "query", "MpsSvc")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrap(err, "error checking windows firewall")
	}
	return strings.Contains(string(output), "RUNNING"), nil
}

func CheckAVExists() (bool, error) {
	type AntiVirusProduct struct {
		DisplayName            string
		ProductState           uint32
		PathToSignedProductExe string
	}

	var antivirus []AntiVirusProduct
	err := wmi.Query("SELECT * FROM AntiVirusProduct", &antivirus)
	if err != nil {
		return false, errors.Wrap(err, "error querying WMI")
	}

	if len(antivirus) > 0 {
		return true, nil
	}

	return false, nil
}
