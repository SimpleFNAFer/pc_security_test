package ui

import (
	"sync"

	"fyne.io/fyne/v2/widget"
)

func newStrListWithImportance() *strListWithImportance {
	return &strListWithImportance{
		data: make(map[string]widget.Importance),
	}
}

type strListWithImportance struct {
	sync.RWMutex
	data map[string]widget.Importance
}

func (swi *strListWithImportance) Set(s string, i widget.Importance) {
	swi.Lock()
	defer swi.Unlock()
	swi.data[s] = i
}
func (swi *strListWithImportance) Get(s string) widget.Importance {
	swi.RLock()
	defer swi.RUnlock()

	i, ok := swi.data[s]
	if !ok {
		return widget.MediumImportance
	}
	return i
}
func (swi *strListWithImportance) Length() int {
	swi.RLock()
	defer swi.RUnlock()
	return len(swi.data)
}
func (swi *strListWithImportance) Clear() {
	swi.Lock()
	defer swi.Unlock()
	for k := range swi.data {
		delete(swi.data, k)
	}
}
