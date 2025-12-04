//go:build unix

package command

import "pc_security_test/tester"

func ProcessFindAVRequest(fAVReq FindAVRequest) {
	found := tester.FindBinariesAndPaths(tester.SourceTypeAV)

	res := FindAVResponse{
		ID:    fAVReq.ID,
		Found: found,
	}

	go AddHistoryEntry(findAVResponseToHistoryEntry(res))
	go func() { findAVResponses <- res }()
}
