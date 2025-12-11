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
	resStr := "тест пройден"
	errStr := "проверка прошла без ошибок"
	if !eRes.Passed {
		resStr =
			"eicar файл не был удалён антивирусной программой за ожидаемое время, проверьте уведомления, " +
				"тест успешно пройден, если антивирус обнаружил файл и вывел предупреждение"
	}
	if eRes.Error != nil {
		errStr = fmt.Sprintf("во время проверки произошла ошибка: %s", eRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s; %s", eRes.ID.String(), resStr, errStr),
	}
}

func AwaitEICARResponse() EICARResponse {
	return <-eicarResponses
}
