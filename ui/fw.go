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

type findFWBlock struct {
	resultOutput *canvas.Text
	checkButton  *widget.Button
}

func newFindFWBlock() *findFWBlock {
	block := &findFWBlock{}

	block.resultOutput = canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	block.resultOutput.Alignment = fyne.TextAlignCenter

	block.checkButton = widget.NewButton("Проверить наличие межсетевого экрана", block.onFWCheckButtonClick)

	return block
}

func (b *findFWBlock) getContainer() *fyne.Container {
	connTesting := container.NewGridWithColumns(
		2,
		b.checkButton,
		b.resultOutput,
	)

	return connTesting
}
func (b *findFWBlock) onFWCheckButtonClick() {
	go command.AddToQueue(command.FindFWRequest{
		ID: uuid.New(),
	})

	go b.awaitAndUpdateUI()

}

func (b *findFWBlock) awaitAndUpdateUI() {
	res := command.AwaitFindFWResponse()

	fyne.Do(func() {
		if res.Error != nil {
			b.resultOutput.Text = res.Error.Error()
			b.resultOutput.Color = theme.Color(theme.ColorNameError)
		} else if len(res.Found) > 0 {
			b.resultOutput.Text = "Межсетевой экран найден"
			b.resultOutput.Color = theme.Color(theme.ColorNameSuccess)
		} else {
			b.resultOutput.Text = "Межсетевой экран не найден"
			b.resultOutput.Color = theme.Color(theme.ColorNameWarning)
		}
		b.resultOutput.Refresh()
	})
}

type testFWBlock struct {
	hostInput    *widget.Entry
	checkButton  *widget.Button
	resultOutput *canvas.Text
}

func newTestFWBlock() *testFWBlock {
	block := &testFWBlock{}

	hostInput := widget.NewEntry()
	hostInput.SetPlaceHolder(defaultHostInputText)

	checkButton := widget.NewButton("Проверить работу межсетевого экрана", block.onPingButtonClick)

	resultOutput := canvas.NewText(defaultResultOutputText, theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	block.hostInput = hostInput
	block.checkButton = checkButton
	block.resultOutput = resultOutput

	return block
}

func (b *testFWBlock) getContainer() *fyne.Container {
	return container.NewGridWithColumns(
		3,
		b.hostInput,
		b.checkButton,
		b.resultOutput,
	)
}

func (b *testFWBlock) onPingButtonClick() {
	host := b.hostInput.Text
	if host == "" {
		host = defaultHostInputText
	}

	go command.AddToQueue(command.TestFWRequest{
		ID:   uuid.New(),
		Host: host,
	})

	go b.awaitAndUpdateUI()
}

func (b *testFWBlock) awaitAndUpdateUI() {
	pRes := command.AwaitTestFWResponse()

	fyne.Do(func() {
		if pRes.Error != nil {
			b.resultOutput.Text = pRes.Error.Error()
			b.resultOutput.Color = theme.Color(theme.ColorNameError)
		} else if !pRes.Available {
			b.resultOutput.Text = fmt.Sprintf("МЭ активен, нет доступа к %s", pRes.Host)
			b.resultOutput.Color = theme.Color(theme.ColorNameSuccess)
		} else {
			b.resultOutput.Text = fmt.Sprintf("МЭ не активен, есть доступ к %s", pRes.Host)
			b.resultOutput.Color = theme.Color(theme.ColorNameWarning)
		}

		b.resultOutput.Refresh()
	})
}
