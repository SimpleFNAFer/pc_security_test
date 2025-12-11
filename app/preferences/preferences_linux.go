//go:build linux

package preferences

import "fyne.io/fyne/v2"

func CheckInitAppPrefs(a fyne.App) {
	checkInitCommonAppPrefs(a)
}

func SetDefaultAll() {
	setDefaultAllCommon()
}

var (
	avBinariesDefVal  = []string{"clamscan", "freshclam", "savd", "esets_daemon"}
	avFilePathsDefVal = []string{}
	fwBinariesDefVal  = []string{"iptables", "ufw", "firewalld", "nft"}
	fwFilePathsDefVal = []string{}
)
