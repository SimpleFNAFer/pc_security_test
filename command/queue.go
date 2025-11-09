package command

import "pc_security_test/config"

type request any

var workerNum = config.NewInt("queue.worker_num", 1)

var queue = make(chan request)

func AddToQueue(command request) {
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
	maxNum := workerNum.Get()
	for range maxNum {
		go queueWorker()
	}
}

func StopQueue() {
	close(queue)
}
