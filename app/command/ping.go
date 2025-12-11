package command

import (
	"fmt"
	"pc_security_test/tester"
	"time"

	"github.com/google/uuid"
)

type PingRequest struct {
	ID   uuid.UUID
	Host string
}
type PingResponse struct {
	ID        uuid.UUID
	Host      string
	Available bool
	Error     error
}

var pingResponses = make(chan PingResponse)

func ProcessPingRequest(pReq PingRequest) {
	available, err := tester.Ping(pReq.Host)

	pRes := PingResponse{
		ID:        pReq.ID,
		Host:      pReq.Host,
		Available: available,
		Error:     err,
	}

	go AddHistoryEntry(pingResponseToHistoryEntry(pRes))
	go func() { pingResponses <- pRes }()
}

func pingRequestToHistoryEntry(pReq PingRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка наличия сетевого подключения, хост: %s", pReq.ID, pReq.Host),
	}
}

func pingResponseToHistoryEntry(pRes PingResponse) Entry {
	availableStr := fmt.Sprintf("есть доступ к %s", pRes.Host)
	errStr := "проверка прошла без ошибок"
	if !pRes.Available {
		availableStr = fmt.Sprintf("нет доступа к %s", pRes.Host)
	}
	if pRes.Error != nil {
		errStr = fmt.Sprintf("во время проверки произошла ошибка: %s", pRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s; %s", pRes.ID.String(), availableStr, errStr),
	}
}

func AwaitPingResponse() PingResponse {
	return <-pingResponses
}
