package ui

import (
	"fmt"
	"pc_security_test/command"
	"pc_security_test/preferences"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const (
	pingTitle = "## Проверка наличия сетевого подключения"
	pingDesc  = `Данный раздел предназначен для проверки наличия сетевого подключения путём обращения к ресурсу.

Ресурс по умолчанию можно отредактировать в настройках в разделе "Сеть".

Для начала проверки нажмите "Проверить". Результат отобразится под полем ввода, а также будет сохранён в истории.
	`
)

type pingBlock struct {
	hostInput    *widget.Entry
	checkButton  *widget.Button
	resultOutput *widget.Label
}

func newPingBlock() *pingBlock {
	block := &pingBlock{}

	hostInput := widget.NewEntry()
	preferences.PingDefaultHost.AddListener(binding.NewDataListener(func() {
		ph, _ := preferences.PingDefaultHost.Get()
		hostInput.SetPlaceHolder(ph)
	}))

	checkButton := widget.NewButton("Проверить", block.onPingButtonClick)

	resultOutput := widget.NewLabel("")
	resultOutput.Wrapping = fyne.TextWrapWord
	resultOutput.Hide()

	block.hostInput = hostInput
	block.checkButton = checkButton
	block.resultOutput = resultOutput

	return block
}

func (b *pingBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(pingDesc)
	desc.Wrapping = fyne.TextWrapWord
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(pingTitle), desc,
		container.NewBorder(nil, nil, nil, b.checkButton, b.hostInput),
		b.resultOutput,
	)
}

func (b *pingBlock) onPingButtonClick() {
	host := b.hostInput.Text
	if host == "" {
		host = b.hostInput.PlaceHolder
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
			b.resultOutput.Importance = widget.DangerImportance
		case pRes.Available:
			b.resultOutput.Text = fmt.Sprintf("Доступ к %s", pRes.Host)
			b.resultOutput.Importance = widget.SuccessImportance
		default:
			b.resultOutput.Text = fmt.Sprintf("Нет доступа к %s", pRes.Host)
			b.resultOutput.Importance = widget.WarningImportance
		}
		b.resultOutput.Refresh()
		b.resultOutput.Show()
	})
}
