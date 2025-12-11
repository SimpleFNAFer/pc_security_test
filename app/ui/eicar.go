package ui

import (
	"pc_security_test/command"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const (
	eicarTitle = "## Проверка работы антивируса"
	//nolint:lll
	eicarDesc = `Данный раздел предназначен для проверки работы антивируса в режиме реального времени с помощью EICAR-теста.

Максимальное количество одновременных тестов, а также длительность ожидания одного теста можно отредактировать в настройках в разделе "Антивирус".

Для начала проверки нажмите "Тест EICAR". Результат отобразится ниже, а также будет сохранён в истории.
	`
)

type eicarBlock struct {
	resultOutput *widget.Label
	checkButton  *widget.Button
}

func newEICARBlock() *eicarBlock {
	block := &eicarBlock{}
	block.resultOutput = widget.NewLabel("")
	block.checkButton = widget.NewButton("Тест EICAR", block.onEICARButtonClick)
	return block
}

func (e *eicarBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(eicarDesc)
	desc.Wrapping = fyne.TextWrapWord
	connTesting := container.NewVBox(
		widget.NewRichTextFromMarkdown(eicarTitle), desc,
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
			e.resultOutput.Importance = widget.DangerImportance
		case eRes.Passed:
			e.resultOutput.Text = "EICAR-тест пройден успешно. Антивирус работает"
			e.resultOutput.Importance = widget.SuccessImportance
		default:
			e.resultOutput.Text = "Проверьте уведомления. Тест успешно пройден, если антивирус вывел предупреждение"
			e.resultOutput.Importance = widget.WarningImportance
		}
		e.resultOutput.Refresh()
	})
}
