package ui

import (
	"fmt"
	"pc_security_test/preferences"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	preferencesWindowTitle = "Настройки"
	commonTitle            = "## Общие"
	appearanceThemeTitle   = "Тема"
	queueWorkerNumTitle    = "Количество параллельно запущенных обработчиков"
	eicarMaxParallelTitle  = "Маскимальное количество одновременных EICAR-тестов"
	netTitle               = "## Сеть"
	defaultHostTitle       = "Адрес по умолчанию"
	avTitle                = "## Антивирус"
	fwTitle                = "## Межсетевой экран"
	binariesTitle          = "Бинарные файлы"
	filepathsTitle         = "Файловые пути"
)

func OpenPreferencesWindow() {
	pw := fyne.CurrentApp().NewWindow(preferencesWindowTitle)
	pw.Resize(fyne.NewSize(400, 800))

	appearanceTheme := widget.NewSelectWithData(preferences.AvailableAppearanceTheme(), preferences.AppearanceTheme)
	appearanceThemeBlock := container.NewHBox(
		widget.NewLabel(appearanceThemeTitle),
		appearanceTheme,
	)

	queueWorkerNum := widget.NewSliderWithData(
		float64(preferences.QueueWorkerNumMin),
		float64(preferences.QueueWorkerNumMax),
		binding.IntToFloat(preferences.QueueWorkerNum),
	)
	queueWorkerNum.Step = 1
	queueWorkerNumBlock := container.NewVBox(
		widget.NewLabel(queueWorkerNumTitle),
		queueWorkerNum,
		sliderLegend(preferences.QueueWorkerNumMin, preferences.QueueWorkerNumMax),
	)

	pingDefaultHost := widget.NewEntryWithData(preferences.PingDefaultHost)
	pingBlock := container.NewBorder(
		widget.NewLabel(defaultHostTitle),
		pingDefaultHost,
		nil, nil, nil,
	)

	eicarMaxParallel := widget.NewSliderWithData(
		float64(preferences.EICARMaxParallelMin),
		float64(preferences.EICARMaxParallelMax),
		binding.IntToFloat(preferences.EICARMaxParallel),
	)
	eicarMaxParallel.Step = 1
	eicarMaxParallelBlock := container.NewVBox(
		widget.NewLabel(eicarMaxParallelTitle),
		eicarMaxParallel,
		sliderLegend(preferences.EICARMaxParallelMin, preferences.EICARMaxParallelMax),
	)

	avBinaries := stringListBlock(preferences.AVBinaries)
	avFilePaths := stringListBlock(preferences.AVFilePaths)
	avTabs := container.NewAppTabs(
		container.NewTabItem(binariesTitle, avBinaries),
		container.NewTabItem(filepathsTitle, avFilePaths),
	)
	avBlock := container.NewVBox(
		widget.NewRichTextFromMarkdown(avTitle),
		eicarMaxParallelBlock,
		widget.NewSeparator(),
		avTabs,
		widget.NewSeparator(),
	)

	fwBinaries := stringListBlock(preferences.FWBinaries)
	fwFilePaths := stringListBlock(preferences.FWFilePaths)
	fwTabs := container.NewAppTabs(
		container.NewTabItem(binariesTitle, fwBinaries),
		container.NewTabItem(filepathsTitle, fwFilePaths),
	)
	fwBlock := container.NewVBox(
		widget.NewRichTextFromMarkdown(fwTitle),
		fwTabs,
		widget.NewSeparator(),
	)

	scroll := container.NewVScroll(
		container.New(
			layout.CustomPaddedLayout{
				LeftPadding:  50,
				RightPadding: 50,
			},
			container.NewVBox(
				widget.NewRichTextFromMarkdown(commonTitle),
				appearanceThemeBlock,
				widget.NewSeparator(),
				queueWorkerNumBlock,
				widget.NewSeparator(),
				widget.NewRichTextFromMarkdown(netTitle),
				pingBlock,
				widget.NewSeparator(),
				avBlock,
				fwBlock,
			),
		),
	)
	scroll.SetMinSize(fyne.NewSize(400, 800))

	pw.CenterOnScreen()
	pw.SetContent(scroll)
	pw.Show()
}

func sliderLegend(min, max int) *fyne.Container {
	objects := []fyne.CanvasObject{}
	for i := min; i < max; i++ {
		l := widget.NewLabel(fmt.Sprintf("%d", i))
		objects = append(objects, l, layout.NewSpacer())
	}
	l := widget.NewLabel(fmt.Sprintf("%d", max))
	objects = append(objects, l)

	return container.NewHBox(objects...)
}

func stringListBlock(strList binding.StringList) *fyne.Container {
	list := widget.NewListWithData(
		strList,
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
		},
	)
	entry := widget.NewEntry()
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		text := entry.Text
		strSlice, _ := strList.Get()
		if text != "" && !slices.Contains(strSlice, text) {
			fyne.LogError("error appending to the list", strList.Append(text))
		}
		entry.SetText("")
	})
	rmBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		text := entry.Text
		strSlice, _ := strList.Get()
		if text != "" && slices.Contains(strSlice, text) {
			fyne.LogError("error removing from the list", strList.Remove(text))
		}
		entry.SetText("")
		list.UnselectAll()
	})
	rmBtn.Disable()
	panel := container.NewBorder(
		nil, nil, nil, container.NewHBox(addBtn, rmBtn), entry,
	)
	list.OnUnselected = func(id widget.ListItemID) { rmBtn.Disable() }
	list.OnSelected = func(id widget.ListItemID) {
		selectedItemText, _ := strList.GetValue(id)
		entry.SetText(selectedItemText)
		rmBtn.Enable()
	}
	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 160))
	block := container.NewBorder(
		panel, nil, nil, nil, scroll,
	)

	return block
}
