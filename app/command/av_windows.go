//go:build windows

package command

import "pc_security_test/tester"

func ProcessFindAVRequest(fAVReq FindAVRequest) {
	found, err := tester.FindAV()

	res := FindAVResponse{
		ID:    fAVReq.ID,
		Found: found,
		Error: err,
	}

	go AddHistoryEntry(findAVResponseToHistoryEntry(res))
	go func() { findAVResponses <- res }()
}
