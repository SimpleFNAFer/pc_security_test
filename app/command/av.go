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

var findAVResponses = make(chan SearchResponse)

func findAVRequestToHistoryEntry(fAVReq FindAVRequest) Entry {
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tПроверка наличия антивируса", fAVReq.ID),
	}
}

func findAVResponseToHistoryEntry(fAVRes SearchResponse) Entry {
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

	errStr := "проверка прошла без ошибок"
	if fAVRes.Error != nil {
		errStr = fmt.Sprintf("во время проверки произошла ошибка: %s", fAVRes.Error.Error())
	}
	return Entry{
		Timestamp: time.Now(),
		Value:     fmt.Sprintf("%s\t|\tРезультат: %s, %s", fAVRes.ID.String(), avsStr, errStr),
	}
}

func AwaitFindAVResponse() SearchResponse {
	return <-findAVResponses
}
