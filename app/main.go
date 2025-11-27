package main

import (
	"pc_security_test/command"
	"pc_security_test/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	command.ProcessQueue()

	ui.InitMasterWindow(a)
	a.Run()

	command.StopQueue()
}
