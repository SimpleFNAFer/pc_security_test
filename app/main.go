package main

import (
	"pc_security_test/command"
	"pc_security_test/preferences"
	"pc_security_test/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("com.sf.pcsecuritycheck")
	preferences.CheckInitAppPrefs(a)

	command.ProcessQueue()
	defer command.StopQueue()

	ui.InitMasterWindow(a)
	a.Run()
}
