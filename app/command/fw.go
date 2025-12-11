package command

import (
	"fmt"
	"pc_security_test/preferences"
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

	errStr := "проверка прошла без ошибок"
	if fFWRes.Error != nil {
		errStr = fmt.Sprintf("во время проверки произошла ошибка: %s", fFWRes.Error.Error())
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
	ID       uuid.UUID
	Protocol preferences.Protocol
	Host     string
	Port     string
}
type TestFWResponse struct {
	ID          uuid.UUID
	Protocol    preferences.Protocol
	Host        string
	Port        string
	Unavailable bool
	Error       error
}

var testFWResponses = make(chan TestFWResponse)

func ProcessTestFWRequest(req TestFWRequest) {
	unavailable, err := tester.FWTest(req.Host, req.Port, req.Protocol)

	res := TestFWResponse{
		ID:          req.ID,
		Protocol:    req.Protocol,
		Host:        req.Host,
		Port:        req.Port,
		Unavailable: unavailable,
		Error:       err,
	}

	go AddHistoryEntry(testFWResponseToHistoryEntry(res))
	go func() { testFWResponses <- res }()
}

func testFWRequestToHistoryEntry(req TestFWRequest) Entry {
	var reqStr strings.Builder
	reqStr.WriteString("хост: " + req.Host)
	if req.Port != "" {
		reqStr.WriteString(", порт: " + req.Port)
	}
	reqStr.WriteString(fmt.Sprintf(", протокол: %s", req.Protocol))
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка работы МЭ: %s", req.ID, reqStr.String()),
	}
}

func testFWResponseToHistoryEntry(res TestFWResponse) Entry {
	var resStr strings.Builder
	resStr.WriteString("к хосту " + res.Host)
	if res.Port != "" {
		resStr.WriteString(fmt.Sprintf(" (порт: %s)", res.Port))
	}
	resStr.WriteString(fmt.Sprintf(" по протоколу %s", res.Protocol))

	availableStr := fmt.Sprintf("МЭ работает некорректно, есть доступ %s", resStr.String())
	errStr := "проверка прошла без ошибок"
	if res.Unavailable {
		availableStr = fmt.Sprintf("МЭ работает корректно, нет доступа %s", res.Host)
	}
	if res.Error != nil {
		errStr = fmt.Sprintf("во время проверки произошла ошибка: %s", res.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s, %s", res.ID.String(), availableStr, errStr),
	}
}

func AwaitTestFWResponse() TestFWResponse {
	return <-testFWResponses
}
