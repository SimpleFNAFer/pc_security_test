package command

import (
	"fmt"
	"pc_security_test/tester"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FindFWRequest struct {
	ID uuid.UUID
}

var findFWResponses = make(chan SearchResponse)

func findFWRequestToHistoryEntry(fFWReq FindFWRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка наличия межсетевого экрана", fFWReq.ID),
	}
}

func findFWResponseToHistoryEntry(fFWRes SearchResponse) Entry {
	fws := []string{}
	for key, value := range fFWRes.Found {
		if value == "" {
			fws = append(fws, key)
		} else {
			fws = append(fws, fmt.Sprintf("%s (%s)", value, key))
		}
	}

	fwsStr := "межсетевой экран не найден"
	if len(fws) > 0 {
		fwsStr = fmt.Sprintf("межсетевой(-ые) экран(-ы): %s", strings.Join(fws, "; "))
	}

	errStr := "успешное выполнение проверки"
	if fFWRes.Error != nil {
		errStr = fmt.Sprintf("ошибка: %s", fFWRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s, %s", fFWRes.ID.String(), fwsStr, errStr),
	}
}

func AwaitFindFWResponse() SearchResponse {
	return <-findFWResponses
}

type TestFWRequest struct {
	ID   uuid.UUID
	Host string
}
type TestFWResponse struct {
	ID        uuid.UUID
	Host      string
	Available bool
	Error     error
}

var testFWResponses = make(chan TestFWResponse)

func ProcessTestFWRequest(req TestFWRequest) {
	available, err := tester.Ping(req.Host)

	res := TestFWResponse{
		ID:        req.ID,
		Host:      req.Host,
		Available: available,
		Error:     err,
	}

	go AddHistoryEntry(testFWResponseToHistoryEntry(res))
	go func() { testFWResponses <- res }()
}

func testFWRequestToHistoryEntry(req TestFWRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка работоспособности МЭ, хост: %s", req.ID, req.Host),
	}
}

func testFWResponseToHistoryEntry(res TestFWResponse) Entry {
	availableStr := fmt.Sprintf("есть доступ к %s", res.Host)
	errStr := "успешное выполнение проверки"
	if !res.Available {
		availableStr = fmt.Sprintf("МЭ работает корректно, нет доступа к %s", res.Host)
	}
	if res.Error != nil {
		errStr = fmt.Sprintf("ошибка: %s", res.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s, %s", res.ID.String(), availableStr, errStr),
	}
}

func AwaitTestFWResponse() TestFWResponse {
	return <-testFWResponses
}
