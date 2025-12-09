package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/gosthome/icons"
	"github.com/gosthome/icons/fynico"
	_ "github.com/gosthome/icons/fynico/google/materialdesignicons"
)

func makeIcon(name string) fyne.Resource {
	icon, err := icons.Parse(fmt.Sprintf("materialdesignicons:%s", name))
	if err != nil {
		fyne.LogError("error parsing icon", err)
	}
	iconRes := fynico.Collections.Lookup(icon.Collection, icon.Icon)
	if iconRes == nil {
		fyne.LogError(fmt.Sprintf("icon %s not found", name), nil)
	}
	return theme.NewColoredResource(iconRes, theme.ColorNameForeground)
}

var (
	lightModeIcon = makeIcon("light_mode")
	darkModeIcon  = makeIcon("dark_mode")
	monitorIcon   = makeIcon("monitor")
)
