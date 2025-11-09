package command

import (
	"fmt"
	"pc_security_test/tester"
	"time"

	"github.com/google/uuid"
)

type EICARRequest struct {
	ID uuid.UUID
}
type EICARResponse struct {
	ID     uuid.UUID
	Passed bool
	Error  error
}

var eicarResponses = make(chan EICARResponse)

func ProcessEICARRequest(eReq EICARRequest) {
	passed, err := tester.EICARTest()

	res := EICARResponse{
		ID:     eReq.ID,
		Passed: passed,
		Error:  err,
	}

	go AddHistoryEntry(eicarResponseToHistoryEntry(res))
	go func() { eicarResponses <- res }()
}

func eicarRequestToHistoryEntry(eReq EICARRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tЗапуск EICAR теста", eReq.ID),
	}
}

func eicarResponseToHistoryEntry(eRes EICARResponse) Entry {
	passStr := "тест пройден"
	errStr := "успешное выполнение"
	if !eRes.Passed {
		passStr = "тест не пройден"
	}
	if eRes.Error != nil {
		errStr = fmt.Sprintf("ошибка: %s", eRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s; %s", eRes.ID.String(), passStr, errStr),
	}
}

func AwaitEICARResponse() EICARResponse {
	return <-eicarResponses
}
