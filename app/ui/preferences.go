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
	eicarWaitDurationTitle = "Время ожидания удаления EICAR-файла"
	netTitle               = "## Сеть"
	defaultHostTitle       = "Адрес по умолчанию"
	avTitle                = "## Антивирус"
	fwTitle                = "## Межсетевой экран"
	binariesTitle          = "Бинарные файлы"
	filepathsTitle         = "Файловые пути"
)

type strPrefsEntry struct {
	widget.Entry
	data binding.String
}

func NewStrPrefsEntry(data binding.String) *strPrefsEntry {
	entry := &strPrefsEntry{
		data: data,
	}
	entry.ExtendBaseWidget(entry)
	entry.OnChanged = entry.onChanged
	entry.OnSubmitted = entry.onSubmitted
	text, _ := data.Get()
	entry.SetText(text)
	return entry
}

func (e *strPrefsEntry) onChanged(s string) {
	if err := e.Validate(); err == nil {
		_ = e.data.Set(s)
	}
}
func (e *strPrefsEntry) FocusLost() {
	e.Entry.FocusLost()
	if err := e.Validate(); err != nil {
		validV, _ := e.data.Get()
		e.SetText(validV)
	}
}
func (e *strPrefsEntry) onSubmitted(s string) {
	if err := e.Validate(); err != nil {
		validV, _ := e.data.Get()
		e.SetText(validV)
	}
}

func OpenPreferencesWindow() {
	pw := fyne.CurrentApp().NewWindow(preferencesWindowTitle)
	pw.Resize(fyne.NewSize(400, 800))

	appearanceTheme := widget.NewSelectWithData(preferences.AvailableAppearanceTheme(), preferences.AppearanceTheme)
	appearanceThemeBlock := container.NewVBox(
		prefTitleWithResetBtn(appearanceThemeTitle, preferences.SetDefaultAppearanceTheme),
		appearanceTheme,
	)

	queueWorkerNum := widget.NewSliderWithData(
		float64(preferences.QueueWorkerNumMin),
		float64(preferences.QueueWorkerNumMax),
		binding.IntToFloat(preferences.QueueWorkerNum),
	)
	queueWorkerNum.Step = 1
	queueWorkerNumBlock := container.NewVBox(
		prefTitleWithResetBtn(queueWorkerNumTitle, preferences.SetDefaultQueueWorkerNum),
		sliderValue(preferences.QueueWorkerNum),
		queueWorkerNum,
		sliderLegend(preferences.QueueWorkerNumMin, preferences.QueueWorkerNumMax),
	)

	pingDefaultHost := NewStrPrefsEntry(preferences.PingDefaultHost)
	pingDefaultHost.Validator = preferences.PingDefaultHostValidator
	pingBlock := container.NewBorder(
		prefTitleWithResetBtn(defaultHostTitle, preferences.SetDefaultPingDefaultHost),
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
		prefTitleWithResetBtn(eicarMaxParallelTitle, preferences.SetDefaultEICARMaxParallel),
		sliderValue(preferences.EICARMaxParallel),
		eicarMaxParallel,
		sliderLegend(preferences.EICARMaxParallelMin, preferences.EICARMaxParallelMax),
	)

	eicarWaitDuration := NewStrPrefsEntry(preferences.EICARWaitDuration)
	eicarWaitDuration.Validator = preferences.EICARWaitDurationValidator
	eicarWaitDurationBlock := container.NewVBox(
		prefTitleWithResetBtn(
			fmt.Sprintf("%s (от %s до %s включительно)",
				eicarWaitDurationTitle,
				preferences.EICARWaitDurationMin,
				preferences.EICARWaitDurationMax,
			),
			preferences.SetDefaultEICARMaxParallel,
		),
		eicarWaitDuration,
	)

	avBinaries := stringListBlock(preferences.AVBinaries, preferences.SetDefaultAVBinaries)
	avFilePaths := stringListBlock(preferences.AVFilePaths, preferences.SetDefaultAVFilePaths)
	avTabs := container.NewAppTabs(
		container.NewTabItem(binariesTitle, avBinaries),
		container.NewTabItem(filepathsTitle, avFilePaths),
	)
	avBlock := container.NewVBox(
		widget.NewRichTextFromMarkdown(avTitle),
		eicarMaxParallelBlock,
		widget.NewSeparator(),
		eicarWaitDurationBlock,
		widget.NewSeparator(),
		avTabs,
		widget.NewSeparator(),
	)

	fwBinaries := stringListBlock(preferences.FWBinaries, preferences.SetDefaultFWBinaries)
	fwFilePaths := stringListBlock(preferences.FWFilePaths, preferences.SetDefaultFWFilePaths)
	fwTabs := container.NewAppTabs(
		container.NewTabItem(binariesTitle, fwBinaries),
		container.NewTabItem(filepathsTitle, fwFilePaths),
	)
	fwBlock := container.NewVBox(
		widget.NewRichTextFromMarkdown(fwTitle),
		fwTabs,
		widget.NewSeparator(),
	)

	restoreAllBtn := widget.NewButtonWithIcon("Сбросить все", theme.ContentUndoIcon(), preferences.SetDefaultAll)
	restoreAllBtn.Importance = widget.DangerImportance

	scroll := container.NewVScroll(
		container.New(
			layout.CustomPaddedLayout{
				LeftPadding:  50,
				RightPadding: 50,
			},
			container.NewVBox(
				container.NewHBox(layout.NewSpacer(), restoreAllBtn),
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
	co := func(v int) *widget.Label {
		l := widget.NewLabel(fmt.Sprintf("%d", v))
		l.TextStyle.Bold = true
		l.SizeName = theme.SizeNameCaptionText
		return l
	}

	return container.NewHBox(co(min), layout.NewSpacer(), co(max))
}

func sliderValue(b binding.Int) *fyne.Container {
	e := widget.NewLabelWithData(binding.IntToString(b))
	e.Importance = widget.HighImportance
	e.TextStyle.Bold = true
	e.SizeName = theme.SizeNameHeadingText
	return container.NewHBox(e, layout.NewSpacer())
}

func stringListBlock(strList binding.StringList, reset func()) *fyne.Container {
	list := widget.NewListWithData(
		strList,
		func() fyne.CanvasObject {
			w := widget.NewLabel("")
			return w
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
		},
	)

	entry := widget.NewEntry()
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		text := entry.Text
		strSlice, _ := strList.Get()
		if text != "" && !slices.Contains(strSlice, text) {
			_ = strList.Append(text)
		}
		entry.SetText("")
	})
	rmBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		text := entry.Text
		strSlice, _ := strList.Get()
		if text != "" && slices.Contains(strSlice, text) {
			_ = strList.Remove(text)
		}
		entry.SetText("")
		list.UnselectAll()
	})
	rmBtn.Disable()
	panel := container.NewBorder(
		nil, nil, nil, container.NewHBox(addBtn, rmBtn, resetBtn(reset)), entry,
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

func resetBtn(reset func()) *widget.Button {
	btn := widget.NewButtonWithIcon("", theme.ContentUndoIcon(), reset)
	btn.Importance = widget.LowImportance
	return btn
}

func prefTitleWithResetBtn(title string, reset func()) *fyne.Container {
	return container.NewBorder(nil, nil, resetBtn(reset), nil, widget.NewLabel(title))
}
