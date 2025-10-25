package ui

import (
	"fmt"
	"pc_security_test/tester"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func initPingBlock() *fyne.Container {
	hostInput := widget.NewEntry()
	hostInput.SetPlaceHolder("mail.ru")

	resultOutput := canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	checkButton := widget.NewButton("Проверить подключение", func() {
		onPingButtonClick(hostInput, resultOutput)
	})

	connTesting := container.NewGridWithColumns(
		//layout.NewSpacer(),
		3,
		hostInput,
		checkButton,
		resultOutput,
		//layout.NewSpacer(),
	)

	return connTesting
}

func onPingButtonClick(e *widget.Entry, c *canvas.Text) {
	host := e.Text
	if host == "" {
		host = "mail.ru"
	}
	avail, err := tester.Ping(host)
	if err != nil {
		c.Text = err.Error()
		c.Color = theme.Color(theme.ColorNameError)
	} else if avail {
		c.Text = fmt.Sprintf("Доступ к %s", host)
		c.Color = theme.Color(theme.ColorNameSuccess)
	} else {
		c.Text = fmt.Sprintf("Нет доступа к %s", host)
		c.Color = theme.Color(theme.ColorNameWarning)
	}
}
