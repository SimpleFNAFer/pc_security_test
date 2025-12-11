//go:build darwin

package preferences

import "fyne.io/fyne/v2"

func CheckInitAppPrefs(a fyne.App) {
	checkInitCommonAppPrefs(a)
}

func SetDefaultAll() {
	setDefaultAllCommon()
}

var (
	avBinariesDefVal  = []string{}
	avFilePathsDefVal = []string{
		"/Applications/NortonSecurity.app",
		"/Applications/AvastSecurity.app",
		"/Applications/Malwarebytes.app",
	}
	fwBinariesDefVal  = []string{"pfctl", "ipfw", "socketfilterfw"}
	fwFilePathsDefVal = []string{}
)
