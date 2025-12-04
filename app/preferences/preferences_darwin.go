//go:build darwin

package preferences

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
