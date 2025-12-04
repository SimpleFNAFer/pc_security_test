//go:build windows

package command

import "pc_security_test/tester"

func ProcessFindFWRequest(fFWReq FindFWRequest) {
	found, err := tester.FindFW()

	fFWRes := FindFWResponse{
		ID:    fFWReq.ID,
		Found: found,
		Error: err,
	}

	go AddHistoryEntry(findFWResponseToHistoryEntry(fFWRes))
	go func() { findFWResponses <- fFWRes }()
}
