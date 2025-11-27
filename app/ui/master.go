package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func InitMasterWindow(app fyne.App) {
	mw := app.NewWindow("Проверка защищённости ПК")
	mw.Resize(fyne.NewSize(800, 600))
	mw.SetFixedSize(true)
	mw.SetMaster()
	mw.Show()

	pingBlock := newPingBlock()
	findFWBlock := newFindFWBlock()
	testFWBlock := newTestFWBlock()
	findAVBlock := newFindAVBlock()
	eicarBlock := newEICARBlock()

	historyBlock := initHistoryBlock(mw)

	mc := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		pingBlock.getContainer(),
		findFWBlock.getContainer(),
		testFWBlock.getContainer(),
		findAVBlock.getContainer(),
		eicarBlock.getContainer(),
		layout.NewSpacer(),
		historyBlock.getContainer(),
	)
	mw.SetContent(mc)
}
