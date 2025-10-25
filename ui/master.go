package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func InitMasterWindow(app fyne.App) {
	mw := app.NewWindow("Проверка защищённости ПК")
	mw.Resize(fyne.NewSize(800, 600))
	mw.SetMaster()

	pingBlock := initPingBlock()
	fwCheckBlock := initFWCheckBlock()
	avCheckBlock := initAVCheckBlock()
	eicarBlock := initEICARBlock()

	mc := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		pingBlock,
		fwCheckBlock,
		avCheckBlock,
		eicarBlock,
		layout.NewSpacer(),
	)
	mw.SetContent(mc)

	mw.Show()
}
