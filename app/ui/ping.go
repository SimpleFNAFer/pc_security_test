package ui

import (
	"fmt"
	"pc_security_test/command"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const (
	defaultHostInputText       = "mail.ru"
	defaultPingCheckButtonText = "Проверить подключение"
	defaultResultOutputText    = "Результат"
)

type pingBlock struct {
	hostInput    *widget.Entry
	checkButton  *widget.Button
	resultOutput *canvas.Text
}

func newPingBlock() *pingBlock {
	block := &pingBlock{}

	hostInput := widget.NewEntry()
	hostInput.SetPlaceHolder(defaultHostInputText)

	checkButton := widget.NewButton(defaultPingCheckButtonText, block.onPingButtonClick)

	resultOutput := canvas.NewText(defaultResultOutputText, theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	block.hostInput = hostInput
	block.checkButton = checkButton
	block.resultOutput = resultOutput

	return block
}

func (b *pingBlock) getContainer() *fyne.Container {
	return container.NewVBox(
		b.hostInput,
		b.checkButton,
		b.resultOutput,
	)
}

func (b *pingBlock) onPingButtonClick() {
	host := b.hostInput.Text
	if host == "" {
		host = defaultHostInputText
	}

	go command.AddToQueue(command.PingRequest{
		ID:   uuid.New(),
		Host: host,
	})

	go b.awaitAndUpdateUI()
}

func (b *pingBlock) awaitAndUpdateUI() {
	pRes := command.AwaitPingResponse()

	fyne.Do(func() {
		switch {
		case pRes.Error != nil:
			b.resultOutput.Text = pRes.Error.Error()
			b.resultOutput.Color = theme.Color(theme.ColorNameError)
		case pRes.Available:
			b.resultOutput.Text = fmt.Sprintf("Доступ к %s", pRes.Host)
			b.resultOutput.Color = theme.Color(theme.ColorNameSuccess)
		default:
			b.resultOutput.Text = fmt.Sprintf("Нет доступа к %s", pRes.Host)
			b.resultOutput.Color = theme.Color(theme.ColorNameWarning)
		}
		b.resultOutput.Refresh()
	})
}
