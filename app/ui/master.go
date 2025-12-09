package ui

import (
	"pc_security_test/preferences"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	// testFWBlock := newTestFWBlock()
	searchAVBlock := newSearchAVBlock()
	eicarBlock := newEICARBlock()

	historyBlock := initHistoryBlock(mw)

	ats := container.NewAppTabs(
		container.NewTabItem(
			navPingBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    50,
					BottomPadding: 50,
					LeftPadding:   50,
					RightPadding:  50,
				},
				pingBlock.getContainer(),
			)),
		container.NewTabItem(
			navAVBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    50,
					BottomPadding: 50,
					LeftPadding:   50,
					RightPadding:  50,
				},
				container.NewVBox(
					searchAVBlock.getContainer(),
					widget.NewSeparator(),
					eicarBlock.getContainer(),
				),
			)),
		container.NewTabItem(
			navFWBtnLabel,
			container.New(
				layout.CustomPaddedLayout{
					TopPadding:    50,
					BottomPadding: 50,
					LeftPadding:   50,
					RightPadding:  50,
				},
				container.NewVBox(
					searchFWBlock.getContainer(),
					// layout.NewSpacer(),
					// testFWBlock.getContainer(),
				),
			)),
	)

	toggleThemeBtn := widget.NewButtonWithIcon("", monitorIcon, CT.toggleThemeVariant)
	preferences.AppearanceTheme.AddListener(binding.NewDataListener(func() {
		t, _ := preferences.AppearanceTheme.Get()
		switch t {
		case preferences.AppearanceThemeLight:
			toggleThemeBtn.SetIcon(lightModeIcon)
		case preferences.AppearanceThemeDark:
			toggleThemeBtn.SetIcon(darkModeIcon)
		case preferences.AppearanceThemeSystem:
			toggleThemeBtn.SetIcon(monitorIcon)
		}
	}))

	mc := container.NewBorder(
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.SettingsIcon(), OpenPreferencesWindow),
			layout.NewSpacer(),
			toggleThemeBtn,
		),
		nil, nil, nil,
		container.NewVSplit(container.NewScroll(ats), historyBlock.getContainer()),
	)

	mw.SetContent(mc)
}
