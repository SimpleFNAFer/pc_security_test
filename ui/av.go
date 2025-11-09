package ui

import (
	"pc_security_test/command"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type findAVBlock struct {
	resultOutput *canvas.Text
	checkButton  *widget.Button
}

func newFindAVBlock() *findAVBlock {
	block := &findAVBlock{}

	block.resultOutput = canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	block.resultOutput.Alignment = fyne.TextAlignCenter

	block.checkButton = widget.NewButton("Проверить наличие антивируса", block.onAVCheckButtonClick)

	return block
}

func (b *findAVBlock) getContainer() *fyne.Container {
	connTesting := container.NewGridWithColumns(
		2,
		b.checkButton,
		b.resultOutput,
	)

	return connTesting
}

func (b *findAVBlock) onAVCheckButtonClick() {
	go command.AddToQueue(command.FindAVRequest{
		ID: uuid.New(),
	})

	go b.awaitAndUpdateUI()

}

func (b *findAVBlock) awaitAndUpdateUI() {
	res := command.AwaitFindAVResponse()

	fyne.Do(func() {
		if res.Error != nil {
			b.resultOutput.Text = res.Error.Error()
			b.resultOutput.Color = theme.Color(theme.ColorNameError)
		} else if len(res.Found) > 0 {
			b.resultOutput.Text = "Антивирус найден"
			b.resultOutput.Color = theme.Color(theme.ColorNameSuccess)
		} else {
			b.resultOutput.Text = "Антивирус не найден"
			b.resultOutput.Color = theme.Color(theme.ColorNameWarning)
		}
		b.resultOutput.Refresh()
	})
}
