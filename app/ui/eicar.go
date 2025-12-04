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

const (
	defaultEICARResultOutputText = "Результат"
	defaultEICARCheckButtonText  = "Тест EICAR"
)

type eicarBlock struct {
	resultOutput *canvas.Text
	checkButton  *widget.Button
}

func newEICARBlock() *eicarBlock {
	block := &eicarBlock{}
	block.resultOutput = canvas.NewText(defaultEICARResultOutputText, theme.Color(theme.ColorNameForeground))
	block.resultOutput.Alignment = fyne.TextAlignCenter

	block.checkButton = widget.NewButton(defaultEICARCheckButtonText, block.onEICARButtonClick)

	return block
}

func (e *eicarBlock) getContainer() *fyne.Container {
	connTesting := container.NewVBox(
		e.checkButton,
		e.resultOutput,
	)
	return connTesting
}
func (e *eicarBlock) onEICARButtonClick() {
	go command.AddToQueue(command.EICARRequest{
		ID: uuid.New(),
	})

	go e.awaitAndUpdateUI()
}

func (e *eicarBlock) awaitAndUpdateUI() {
	eRes := command.AwaitEICARResponse()

	fyne.Do(func() {
		switch {
		case eRes.Error != nil:
			e.resultOutput.Text = eRes.Error.Error()
			e.resultOutput.Color = theme.Color(theme.ColorNameError)
		case eRes.Passed:
			e.resultOutput.Text = "EICAR-тест пройден успешно. Антивирус работает"
			e.resultOutput.Color = theme.Color(theme.ColorNameSuccess)
		default:
			e.resultOutput.Text = "EICAR-тест не пройден. Антивирус не работает"
			e.resultOutput.Color = theme.Color(theme.ColorNameWarning)
		}
		e.resultOutput.Refresh()
	})
}
