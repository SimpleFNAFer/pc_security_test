package main

import (
	"pc_security_test/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	ui.InitMasterWindow(a)
	a.Run()
}
