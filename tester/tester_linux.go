//go:build linux
// +build linux

package tester

import (
	"os/exec"
)

func CheckFWExists() (bool, error) {
	fwBinaries := []string{"ufw", "iptables", "nft"}
	for _, fw := range fwBinaries {
		if _, err := exec.LookPath(fw); err == nil {
			return true, nil
		}
	}

	return false, nil
}

func CheckAVExists() (bool, error) {
	avBinaries := []string{"clamscan", "freshclam", "savd", "esets_daemon"}
	for _, av := range avBinaries {
		if _, err := exec.LookPath(av); err == nil {
			return true, nil
		}
	}
	return false, nil
}
