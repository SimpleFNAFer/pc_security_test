package ui

import (
	"pc_security_test/tester"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func initAVCheckBlock() *fyne.Container {
	resultOutput := canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	checkButton := widget.NewButton("Проверить наличие антивируса", func() {
		onAVCheckButtonClick(resultOutput)
	})

	connTesting := container.NewGridWithColumns(
		//layout.NewSpacer(),
		2,
		checkButton,
		resultOutput,
		//layout.NewSpacer(),
	)

	return connTesting
}
func onAVCheckButtonClick(c *canvas.Text) {
	avWorks, err := tester.CheckAVExists()
	if err != nil {
		c.Text = err.Error()
		c.Color = theme.Color(theme.ColorNameError)
	} else if avWorks {
		c.Text = "Антивирус найден"
		c.Color = theme.Color(theme.ColorNameSuccess)
	} else {
		c.Text = "Антивирус не найден"
		c.Color = theme.Color(theme.ColorNameWarning)
	}
}
