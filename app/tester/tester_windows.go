//go:build windows

package tester

import (
	"fmt"
	"maps"
	"os/exec"
	"pc_security_test/preferences"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/yusufpapurcu/wmi"
)

func FindFW() (map[string]string, error) {
	var (
		res = make(map[string]string)
		err error
	)

	sd, err := preferences.FWSearchDefault.Get()
	if err != nil {
		fyne.LogError("FindFW.FWSearchDefault.Get", err)
	}
	if sd {
		var output []byte
		cmd := exec.Command("sc", "query", "MpsSvc")
		output, err = cmd.CombinedOutput()
		if strings.Contains(string(output), "RUNNING") {
			res["Брандмауэр Windows"] = ""
		}
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

	sd, err := preferences.AVSearchDefault.Get()
	if err != nil {
		fyne.LogError("FindAV.AVSearchDefault.Get", err)
	}
	if sd {
		type AntiVirusProduct struct {
			DisplayName              string
			ProductState             uint32
			PathToSignedProductExe   string
			PathToSignedReportingExe string
		}

		var antivirus []AntiVirusProduct
		err = wmi.QueryNamespace("SELECT * FROM AntivirusProduct", &antivirus, "root\\SecurityCenter2")
		if err != nil {
			err = wmi.QueryNamespace("SELECT * FROM AntivirusProduct", &antivirus, "root\\SecurityCenter")
		}

		for _, av := range antivirus {
			res[av.DisplayName] = fmt.Sprintf("%s ; %s", av.PathToSignedProductExe, av.PathToSignedReportingExe)
		}
	}

	found := FindBinariesAndPaths(SourceTypeAV)
	maps.Copy(res, found)

	return res, err
}
