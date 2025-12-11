//go:build windows

package ui

import (
	"pc_security_test/preferences"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const (
	avSearchDefaultTitle = "Проверять наличие антивируса по умолчанию"
	fwSearchDefaultTitle = "Проверять наличие межсетевого экрана по умолчанию"
)

func avBlock() *fyne.Container {
	b := avBlockCommon()
	b.Add(prefBoolBlock(avSearchDefaultTitle, preferences.AVSearchDefault, preferences.SetDefaultAVSearchDefault))
	return b
}

func fwBlock() *fyne.Container {
	b := fwBlockCommon()
	b.Add(prefBoolBlock(fwSearchDefaultTitle, preferences.FWSearchDefault, preferences.SetDefaultFWSearchDefault))
	return b
}

func prefBoolBlock(title string, b binding.Bool, reset func()) *fyne.Container {
	return container.NewBorder(
		nil, nil, nil, resetBtn(reset),
		widget.NewCheckWithData(title, b),
	)
}
