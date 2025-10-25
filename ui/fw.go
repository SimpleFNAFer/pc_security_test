package ui

import (
	"pc_security_test/tester"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func initFWCheckBlock() *fyne.Container {
	resultOutput := canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	checkButton := widget.NewButton("Проверить наличие межсетевого экрана", func() {
		onFWCheckButtonClick(resultOutput)
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
func onFWCheckButtonClick(c *canvas.Text) {
	avWorks, err := tester.CheckFWExists()
	if err != nil {
		c.Text = err.Error()
		c.Color = theme.Color(theme.ColorNameError)
	} else if avWorks {
		c.Text = "Межсетевой экран найден"
		c.Color = theme.Color(theme.ColorNameSuccess)
	} else {
		c.Text = "Межсетевой экран не найден"
		c.Color = theme.Color(theme.ColorNameWarning)
	}
}
