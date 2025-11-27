//go:build windows
// +build windows

package tester

import (
	"maps"
	"os/exec"
	"pc_security_test/config"
	"strings"

	"github.com/yusufpapurcu/wmi"
)

var (
	findAVBinariesSlice = config.NewStringSlice("find_av.windows.binaries", []string{})
	findAVPathsSlice    = config.NewStringSlice("find_av.windows.paths", []string{})

	findFWBinariesSlice = config.NewStringSlice("find_fw.windows.binaries", []string{})
	findFWPathsSlice    = config.NewStringSlice("find_fw.windows.paths", []string{})
)

func FindFW() (map[string]string, error) {
	var (
		res = make(map[string]string)
		err error
	)

	cmd := exec.Command("sc", "query", "MpsSvc")
	output, err := cmd.CombinedOutput()
	if strings.Contains(string(output), "RUNNING") {
		res["Брандмауэр Windows"] = ""
	}

	found := FindBinariesAndPaths(SourceTypeFW)
	maps.Copy(res, found)

	return res, err
}

func FindAV() (map[string]string, error) {
	var (
		res = make(map[string]string)
		err error
	)

	type AntiVirusProduct struct {
		DisplayName            string
		ProductState           uint32
		PathToSignedProductExe string
	}

	var antivirus []AntiVirusProduct
	err = wmi.Query("SELECT * FROM AntiVirusProduct", &antivirus)

	for _, av := range antivirus {
		res[av.DisplayName] = av.PathToSignedProductExe
	}

	found := FindBinariesAndPaths(SourceTypeAV)
	maps.Copy(res, found)

	return res, err
}
