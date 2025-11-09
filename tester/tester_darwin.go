//go:build darwin
// +build darwin

package tester

import (
	"pc_security_test/config"
)

var (
	findAVBinariesSlice = config.NewStringSlice("find_av.darwin.binaries", []string{})
	findAVPathsSlice    = config.NewStringSlice("find_av.darwin.paths", []string{})

	findFWBinariesSlice = config.NewStringSlice("find_fw.darwin.binaries", []string{})
	findFWPathsSlice    = config.NewStringSlice("find_fw.darwin.paths", []string{})
)
