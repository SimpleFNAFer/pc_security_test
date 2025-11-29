package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	prefsKeyTheme        = "theme"
	prefsValueThemeLight = "light"
	prefsValueThemeDark  = "dark"
)

type customTheme struct {
	fyne.Theme
	isDark bool
}

func initAndSetCustomTheme() *customTheme {
	if fyne.CurrentApp().Preferences().String(prefsKeyTheme) == "" {
		fyne.CurrentApp().Preferences().SetString(prefsKeyTheme, prefsValueThemeLight)
	}

	ct := customTheme{
		isDark: fyne.CurrentApp().Preferences().String(prefsKeyTheme) == prefsValueThemeDark,
	}

	fyne.CurrentApp().Settings().SetTheme(&ct)
	return &ct
}

func (ct *customTheme) Color(tcn fyne.ThemeColorName, tv fyne.ThemeVariant) color.Color {
	if ct.isDark {
		return theme.DefaultTheme().Color(tcn, theme.VariantDark)
	}
	return theme.DefaultTheme().Color(tcn, theme.VariantLight)
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
	ct.isDark = !ct.isDark
	variant := prefsValueThemeLight
	if ct.isDark {
		variant = prefsValueThemeDark
	}
	fyne.CurrentApp().Preferences().SetString(prefsKeyTheme, variant)
	fyne.CurrentApp().Settings().SetTheme(ct)
}
