package command

import (
	"pc_security_test/preferences"

	"fyne.io/fyne/v2"
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
	maxNum, err := preferences.QueueWorkerNum.Get()
	if err != nil {
		fyne.LogError("ProcessQueue.QueueWorkerNum.Get", err)
	}
	for range maxNum {
		go queueWorker()
	}
}

func StopQueue() {
	close(queue)
}
