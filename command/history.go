package command

import (
	"time"
)

type Entry struct {
	Timestamp time.Time
	Value     string
}

var historyEntries = make(chan Entry)

func AddHistoryEntry(e Entry) {
	historyEntries <- e
}

func HistoryEntries() <-chan Entry {
	return historyEntries
}
