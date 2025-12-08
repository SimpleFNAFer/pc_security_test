package command

import (
	"pc_security_test/preferences"

	"github.com/google/uuid"
)

type Request any
type SearchResponse struct {
	ID    uuid.UUID
	Found map[string]string
	Error error
}

var queue = make(chan Request)

func AddToQueue(command Request) {
	queue <- command
}

func queueWorker() {
	for req := range queue {
		switch typedReq := req.(type) {
		case PingRequest:
			AddHistoryEntry(pingRequestToHistoryEntry(typedReq))
			ProcessPingRequest(typedReq)
		case FindFWRequest:
			AddHistoryEntry(findFWRequestToHistoryEntry(typedReq))
			ProcessFindFWRequest(typedReq)
		case FindAVRequest:
			AddHistoryEntry(findAVRequestToHistoryEntry(typedReq))
			ProcessFindAVRequest(typedReq)
		case EICARRequest:
			AddHistoryEntry(eicarRequestToHistoryEntry(typedReq))
			ProcessEICARRequest(typedReq)
		case TestFWRequest:
			AddHistoryEntry(testFWRequestToHistoryEntry(typedReq))
			ProcessTestFWRequest(typedReq)
		}
	}
}

func ProcessQueue() {
	maxNum, _ := preferences.QueueWorkerNum.Get()
	for range maxNum {
		go queueWorker()
	}
}

func StopQueue() {
	close(queue)
}
