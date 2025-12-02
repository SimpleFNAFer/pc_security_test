package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	navPingBtnLabel = "Сеть"
	navAVBtnLabel   = "Антивирус"
	navFWBtnLabel   = "Межсетевой экран"
)

func InitMasterWindow(app fyne.App) {
	mw := app.NewWindow("Проверка защищённости ПК")
	mw.Resize(fyne.NewSize(800, 800))
	mw.SetFixedSize(true)
	mw.SetMaster()
	mw.Show()

	initAndSetCustomTheme()

	pingBlock := newPingBlock()
	findFWBlock := newFindFWBlock()
	testFWBlock := newTestFWBlock()
	findAVBlock := newFindAVBlock()
	eicarBlock := newEICARBlock()

	historyBlock := initHistoryBlock(mw)

	ats := container.NewAppTabs(
		container.NewTabItem(navPingBtnLabel, container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(
				layout.NewSpacer(),
				pingBlock.getContainer(),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		)),
		container.NewTabItem(navAVBtnLabel, container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(
				layout.NewSpacer(),
				findAVBlock.getContainer(),
				layout.NewSpacer(),
				eicarBlock.getContainer(),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		)),
		container.NewTabItem(navFWBtnLabel, container.NewVBox(
			layout.NewSpacer(),
			container.NewHBox(
				layout.NewSpacer(),
				findFWBlock.getContainer(),
				layout.NewSpacer(),
				testFWBlock.getContainer(),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		)),
	)
	mc := container.NewBorder(
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.SettingsIcon(), OpenPreferencesWindow),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), CT.toggleThemeVariant),
		),
		historyBlock.getContainer(), nil, nil,
		ats,
	)

	mw.SetContent(mc)
}
