//go:build unix
// +build unix

package command

import "pc_security_test/tester"

func ProcessFindFWRequest(fFWReq FindFWRequest) {
	found := tester.FindBinariesAndPaths(tester.SourceTypeFW)

	fFWRes := FindFWResponse{
		ID:    fFWReq.ID,
		Found: found,
	}

	go AddHistoryEntry(findFWResponseToHistoryEntry(fFWRes))
	go func() { findFWResponses <- fFWRes }()
}
