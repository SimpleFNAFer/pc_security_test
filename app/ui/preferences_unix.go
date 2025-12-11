//go:build unix

package ui

import (
	"fyne.io/fyne/v2"
)

func avBlock() *fyne.Container {
	return avBlockCommon()
}

func fwBlock() *fyne.Container {
	return fwBlockCommon()
}
