package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FindAVRequest struct {
	ID uuid.UUID
}
type FindAVResponse struct {
	ID    uuid.UUID
	Found map[string]string
	Error error
}

var findAVResponses = make(chan FindAVResponse)

func findAVRequestToHistoryEntry(fAVReq FindAVRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка наличия антивируса", fAVReq.ID),
	}
}

func findAVResponseToHistoryEntry(fAVRes FindAVResponse) Entry {
	avs := []string{}
	for key, value := range fAVRes.Found {
		if value == "" {
			avs = append(avs, key)
		} else {
			avs = append(avs, fmt.Sprintf("%s (%s)", value, key))
		}
	}

	avsStr := "антивирус не найден"
	if len(avs) > 0 {
		avsStr = fmt.Sprintf("антивирус(-ы): %s", strings.Join(avs, "; "))
	}

	errStr := "успешное выполнение проверки"
	if fAVRes.Error != nil {
		errStr = fmt.Sprintf("ошибка: %s", fAVRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s, %s", fAVRes.ID.String(), avsStr, errStr),
	}
}

func AwaitFindAVResponse() FindAVResponse {
	return <-findAVResponses
}
