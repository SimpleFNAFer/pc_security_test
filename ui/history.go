package ui

import (
	"fmt"
	"log"
	"pc_security_test/command"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	saveToFileBtnText   = "Сохранить в файл"
	clearHistoryBtnText = "Очистить историю"
)

type historyBlock struct {
	master fyne.Window

	textLines       *widget.TextGrid
	saveToFileBtn   *widget.Button
	clearHistoryBtn *widget.Button
	sc              *container.Scroll
}

func initHistoryBlock(master fyne.Window) *historyBlock {
	if reflect.ValueOf(master).IsNil() {
		panic("master window is nil")
	}

	block := &historyBlock{
		master:    master,
		textLines: widget.NewTextGrid(),
	}

	block.textLines.ShowLineNumbers = true
	block.saveToFileBtn = widget.NewButton(saveToFileBtnText, block.saveToFile)
	block.clearHistoryBtn = widget.NewButton(clearHistoryBtnText, block.clearHistory)

	sc := container.NewScroll(block.textLines)
	sc.SetMinSize(fyne.NewSize(0, 200))
	block.sc = sc

	return block
}

func (h *historyBlock) getContainer() *fyne.Container {
	buttons := container.NewHBox(
		h.saveToFileBtn,
		h.clearHistoryBtn,
	)

	c := container.NewVBox(h.sc, layout.NewSpacer(), buttons)

	go h.awaitHistoryEntries()

	return c
}

func (h *historyBlock) awaitHistoryEntries() {
	for e := range command.HistoryEntries() {
		timestamp := e.Timestamp.Format("01-02-2006 15:04:05")
		fyne.DoAndWait(func() {
			h.textLines.Append(fmt.Sprintf("%s\t|\t%s", timestamp, e.Value))
			h.textLines.Refresh()
			h.sc.ScrollToBottom()
			h.sc.Refresh()
		})
	}
}

func (h *historyBlock) saveToFile() {
	textBytes := []byte(h.textLines.Text())

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}

		defer writer.Close()

		_, writeErr := writer.Write(textBytes)
		if writeErr != nil {
			log.Println("error writing to file:", writeErr)
			return
		}
	}, h.master)
}

func (h *historyBlock) clearHistory() {
	h.textLines.Rows = nil
	h.textLines.Refresh()
}
