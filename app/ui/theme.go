package ui

import (
	"image/color"
	"pc_security_test/preferences"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct{}

var CT *customTheme

func initAndSetCustomTheme() {
	CT = &customTheme{}
	preferences.AppearanceTheme.AddListener(binding.NewDataListener(func() {
		fyne.CurrentApp().Settings().SetTheme(CT)
	}))
}

func (ct *customTheme) Color(tcn fyne.ThemeColorName, tv fyne.ThemeVariant) color.Color {
	t, err := preferences.AppearanceTheme.Get()
	if err != nil {
		fyne.LogError("Color.AppearanceTheme.Get", err)
	}

	switch t {
	case preferences.AppearanceThemeDark:
		return theme.DefaultTheme().Color(tcn, theme.VariantDark)
	case preferences.AppearanceThemeLight:
		return theme.DefaultTheme().Color(tcn, theme.VariantLight)
	default:
		return theme.DefaultTheme().Color(tcn, tv)
	}
}

func (ct *customTheme) Font(ts fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(ts)
}

func (ct *customTheme) Icon(tin fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(tin)
}

func (ct *customTheme) Size(tsn fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(tsn)
}

func (ct *customTheme) toggleThemeVariant() {
	avail := preferences.AvailableAppearanceTheme()
	t, err := preferences.AppearanceTheme.Get()
	if err != nil {
		fyne.LogError("toggleThemeVariant.AppearanceTheme.Get", err)
	}
	ti := slices.Index(avail, t)
	newTI := (ti - 1 + len(avail)) % len(avail)
	newT := avail[newTI]
	if err := preferences.AppearanceTheme.Set(newT); err != nil {
		fyne.LogError("toggleThemeVariant.AppearanceTheme.Set", err)
	}
}
