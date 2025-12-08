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
	mw.Resize(fyne.NewSize(800, 600))
	mw.SetMaster()
	mw.Show()

	initAndSetCustomTheme()

	pingBlock := newPingBlock()
	searchFWBlock := newSearchFWBlock()
	// todo testFWBlock := newTestFWBlock()
	searchAVBlock := newSearchAVBlock()
	eicarBlock := newEICARBlock()

	historyBlock := initHistoryBlock(mw)

	ats := container.NewAppTabs(
		container.NewTabItem(
			navPingBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:   50,
					LeftPadding:  50,
					RightPadding: 50,
				},
				pingBlock.getContainer(),
			)),
		container.NewTabItem(
			navAVBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:   50,
					LeftPadding:  50,
					RightPadding: 50,
				},
				container.NewVBox(
					searchAVBlock.getContainer(),
					layout.NewSpacer(),
					eicarBlock.getContainer(),
				),
			)),
		container.NewTabItem(
			navFWBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:   50,
					LeftPadding:  50,
					RightPadding: 50,
				},
				container.NewVBox(
					searchFWBlock.getContainer(),
					//layout.NewSpacer(),
					//testFWBlock.getContainer(),
				),
			)),
	)
	mc := container.NewBorder(
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.SettingsIcon(), OpenPreferencesWindow),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), CT.toggleThemeVariant),
		),
		nil, nil, nil,
		container.NewVSplit(container.NewScroll(ats), historyBlock.getContainer()),
	)

	mw.SetContent(mc)
}
