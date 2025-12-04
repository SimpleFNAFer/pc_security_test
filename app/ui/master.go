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
		container.NewTabItem(
			navPingBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    100,
					BottomPadding: 100,
					LeftPadding:   200,
					RightPadding:  200,
				},
				pingBlock.getContainer(),
			)),
		container.NewTabItem(
			navAVBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    100,
					BottomPadding: 100,
					LeftPadding:   200,
					RightPadding:  200,
				},
				container.NewHBox(
					findAVBlock.getContainer(),
					layout.NewSpacer(),
					eicarBlock.getContainer(),
				),
			)),
		container.NewTabItem(
			navFWBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    100,
					BottomPadding: 100,
					LeftPadding:   200,
					RightPadding:  200,
				},
				container.NewHBox(
					findFWBlock.getContainer(),
					layout.NewSpacer(),
					testFWBlock.getContainer(),
				),
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
