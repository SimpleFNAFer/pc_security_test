package ui

import (
	"pc_security_test/tester"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func initEICARBlock() *fyne.Container {
	resultOutput := canvas.NewText("Результат", theme.Color(theme.ColorNameForeground))
	resultOutput.Alignment = fyne.TextAlignCenter

	checkButton := widget.NewButton("Тест EICAR", func() {
		onEICARButtonClick(resultOutput)
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
func onEICARButtonClick(c *canvas.Text) {
	avWorks, err := tester.EICARTest()
	if err != nil {
		c.Text = err.Error()
		c.Color = theme.Color(theme.ColorNameError)
	} else if avWorks {
		c.Text = "EICAR-тест пройден успешно. Антивирус работает"
		c.Color = theme.Color(theme.ColorNameSuccess)
	} else {
		c.Text = "EICAR-тест не пройден. Антивирус не работает"
		c.Color = theme.Color(theme.ColorNameWarning)
	}
}
