//go:build linux
// +build linux

package tester

import (
	"pc_security_test/config"
)

var (
	findAVBinariesSlice = config.NewStringSlice("find_av.linux.binaries", []string{})
	findAVPathsSlice    = config.NewStringSlice("find_av.linux.paths", []string{})

	findFWBinariesSlice = config.NewStringSlice("find_fw.linux.binaries", []string{})
	findFWPathsSlice    = config.NewStringSlice("find_fw.linux.paths", []string{})
)
