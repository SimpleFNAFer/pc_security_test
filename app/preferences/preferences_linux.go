//go:build linux
// +build linux

package preferences

var (
	avBinariesDefVal  = []string{"clamscan", "freshclam", "savd", "esets_daemon"}
	avFilePathsDefVal = []string{}
	fwBinariesDefVal  = []string{"iptables", "ufw", "firewalld", "nft"}
	fwFilePathsDefVal = []string{}
)
