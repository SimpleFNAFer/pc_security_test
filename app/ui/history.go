package ui

import (
	"fmt"
	"pc_security_test/command"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	saveToFileBtnText   = "Сохранить в файл"
	clearHistoryBtnText = "Очистить историю"
)

type historyBlock struct {
	master fyne.Window

	textLines       *widget.Label
	saveToFileBtn   *widget.Button
	clearHistoryBtn *widget.Button
	sc              *container.Scroll
}

func initHistoryBlock(master fyne.Window) *historyBlock {
	block := &historyBlock{
		master:    master,
		textLines: widget.NewLabel(""),
	}

	block.textLines.Wrapping = fyne.TextWrapBreak
	block.saveToFileBtn = widget.NewButton(saveToFileBtnText, block.saveToFile)
	block.clearHistoryBtn = widget.NewButton(clearHistoryBtnText, block.clearHistory)

	sc := container.NewVScroll(block.textLines)
	block.sc = sc

	return block
}

func (h *historyBlock) getContainer() *fyne.Container {
	buttons := container.NewHBox(
		h.saveToFileBtn,
		h.clearHistoryBtn,
	)

	c := container.NewBorder(widget.NewSeparator(), buttons, nil, nil, h.sc)

	go h.awaitHistoryEntries()

	return c
}

func (h *historyBlock) awaitHistoryEntries() {
	for e := range command.HistoryEntries() {
		timestamp := e.Timestamp.Format("01-02-2006 15:04:05")
		fyne.DoAndWait(func() {
			h.textLines.SetText(fmt.Sprintf("%s\n%s\t|\t%s", h.textLines.Text, timestamp, e.Value))
			h.textLines.Refresh()
			h.sc.ScrollToBottom()
			h.sc.Refresh()
		})
	}
}

func (h *historyBlock) saveToFile() {
	textBytes := []byte(h.textLines.Text)

	fileSaveDlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}

		defer func() {
			if err := writer.Close(); err != nil {
				fyne.LogError("error closing writer", err)
			}
		}()

		_, writeErr := writer.Write(textBytes)
		if writeErr != nil {
			fyne.LogError("error writing to file", writeErr)
			return
		}
	}, h.master)

	fileSaveDlg.Resize(fyne.NewSize(800, 600))
	fileSaveDlg.Show()
}

func (h *historyBlock) clearHistory() {
	h.textLines.SetText("")
	h.textLines.Refresh()
}
