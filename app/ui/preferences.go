package ui

import (
	"fmt"
	"pc_security_test/preferences"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	preferencesWindowTitle   = "Настройки"
	commonTitle              = "## Общие"
	appearanceThemeTitle     = "Тема"
	queueWorkerNumTitle      = "Количество параллельно запущенных обработчиков"
	eicarMaxParallelTitle    = "Маскимальное количество одновременных EICAR-тестов"
	eicarWaitDurationTitle   = "Время ожидания удаления EICAR-файла"
	pingWaitDurationTitle    = "Время ожидания icmp-ответа"
	fwCheckWaitDurationTitle = "Время ожидания подключения"
	netTitle                 = "## Сеть"
	defaultHostTitle         = "Адрес по умолчанию"
	avTitle                  = "## Антивирус"
	fwTitle                  = "## Межсетевой экран"
	binariesTitle            = "Бинарные файлы"
	filepathsTitle           = "Файловые пути"
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
	data.AddListener(binding.NewDataListener(func() {
		text, err := data.Get()
		if err != nil {
			fyne.LogError("NewStrPrefsEntry.data.Get", err)
		}
		entry.SetText(text)
	}))
	return entry
}

func (e *strPrefsEntry) onChanged(s string) {
	if err := e.Validate(); err == nil {
		if err := e.data.Set(s); err != nil {
			fyne.LogError("onChanged", err)
		}
	}
}
func (e *strPrefsEntry) FocusLost() {
	e.Entry.FocusLost()
	if err := e.Validate(); err != nil {
		validV, err := e.data.Get()
		if err != nil {
			fyne.LogError("FocusLost.data.Get", err)
		}
		e.SetText(validV)
	}
}
func (e *strPrefsEntry) onSubmitted(s string) {
	if err := e.Validate(); err != nil {
		validV, err := e.data.Get()
		if err != nil {
			fyne.LogError("onSubmitted.data.Get", err)
		}
		e.SetText(validV)
	}
}

func OpenPreferencesWindow() {
	pw := fyne.CurrentApp().NewWindow(preferencesWindowTitle)
	pw.Resize(fyne.NewSize(400, 800))

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
				appearanceThemeBlock(),
				widget.NewSeparator(),
				queueWorkerNumBlock(),
				widget.NewSeparator(),
				widget.NewRichTextFromMarkdown(netTitle),
				pingDefaultHostBlock(),
				pingWaitDurationBlock(),
				widget.NewSeparator(),
				avBlock(),
				fwBlock(),
			),
		),
	)
	scroll.SetMinSize(fyne.NewSize(400, 800))

	pw.CenterOnScreen()
	pw.SetContent(scroll)
	pw.Show()
}

func appearanceThemeBlock() *fyne.Container {
	appearanceTheme := widget.NewSelectWithData(preferences.AvailableAppearanceTheme(), preferences.AppearanceTheme)
	return container.NewVBox(
		prefTitleWithResetBtn(appearanceThemeTitle, preferences.SetDefaultAppearanceTheme),
		appearanceTheme,
	)
}

func queueWorkerNumBlock() *fyne.Container {
	queueWorkerNum := widget.NewSliderWithData(
		float64(preferences.QueueWorkerNumMin),
		float64(preferences.QueueWorkerNumMax),
		binding.IntToFloat(preferences.QueueWorkerNum),
	)
	queueWorkerNum.Step = 1
	return container.NewVBox(
		prefTitleWithResetBtn(
			fmt.Sprintf("%s (изменится после перезапуска)", queueWorkerNumTitle),
			preferences.SetDefaultQueueWorkerNum),
		sliderValue(preferences.QueueWorkerNum),
		queueWorkerNum,
		sliderLegend(preferences.QueueWorkerNumMin, preferences.QueueWorkerNumMax),
	)
}

func pingDefaultHostBlock() *fyne.Container {
	pingDefaultHost := NewStrPrefsEntry(preferences.PingDefaultHost)
	pingDefaultHost.Validator = preferences.HostValidator
	return container.NewBorder(
		prefTitleWithResetBtn(defaultHostTitle, preferences.SetDefaultPingDefaultHost),
		pingDefaultHost,
		nil, nil, nil,
	)
}

func pingWaitDurationBlock() *fyne.Container {
	return durationBlock(
		pingWaitDurationTitle,
		preferences.PingWaitDurationMin,
		preferences.PingWaitDurationMax,
		preferences.PingWaitDuration,
		preferences.PingWaitDurationValidator,
		preferences.SetDefaultPingWaitDuration,
	)
}

func eicarMaxParallelBlock() *fyne.Container {
	eicarMaxParallel := widget.NewSliderWithData(
		float64(preferences.EICARMaxParallelMin),
		float64(preferences.EICARMaxParallelMax),
		binding.IntToFloat(preferences.EICARMaxParallel),
	)
	eicarMaxParallel.Step = 1
	return container.NewVBox(
		prefTitleWithResetBtn(eicarMaxParallelTitle, preferences.SetDefaultEICARMaxParallel),
		sliderValue(preferences.EICARMaxParallel),
		eicarMaxParallel,
		sliderLegend(preferences.EICARMaxParallelMin, preferences.EICARMaxParallelMax),
	)
}

func eicarWaitDurationBlock() *fyne.Container {
	return durationBlock(
		eicarWaitDurationTitle,
		preferences.EICARWaitDurationMin,
		preferences.EICARWaitDurationMax,
		preferences.EICARWaitDuration,
		preferences.EICARWaitDurationValidator,
		preferences.SetDefaultEICARWaitDuration,
	)
}

func fwCheckWaitDurationBlock() *fyne.Container {
	return durationBlock(
		fwCheckWaitDurationTitle,
		preferences.FWCheckWaitDurationMin,
		preferences.FWCheckWaitDurationMax,
		preferences.FWCheckWaitDuration,
		preferences.FWCheckWaitDurationValidator,
		preferences.SetDefaultFWCheckWaitDuration,
	)
}

type tabsParams struct {
	title string
	data  binding.StringList
	reset func()
}

func tabs(tabs ...tabsParams) *container.AppTabs {
	c := container.NewAppTabs()
	for _, tp := range tabs {
		l := stringListBlock(tp.data, tp.reset)
		c.Append(container.NewTabItem(tp.title, l))
	}
	return c
}

func avTabs() *container.AppTabs {
	return tabs(
		tabsParams{title: binariesTitle, data: preferences.AVBinaries, reset: preferences.SetDefaultAVBinaries},
		tabsParams{title: filepathsTitle, data: preferences.AVFilePaths, reset: preferences.SetDefaultAVFilePaths},
	)
}

func fwTabs() *container.AppTabs {
	return tabs(
		tabsParams{title: binariesTitle, data: preferences.FWBinaries, reset: preferences.SetDefaultFWBinaries},
		tabsParams{title: filepathsTitle, data: preferences.FWFilePaths, reset: preferences.SetDefaultFWFilePaths},
	)
}

func avBlockCommon() *fyne.Container {
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(avTitle),
		eicarMaxParallelBlock(),
		widget.NewSeparator(),
		eicarWaitDurationBlock(),
		widget.NewSeparator(),
		avTabs(),
		widget.NewSeparator(),
	)
}

func fwBlockCommon() *fyne.Container {
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(fwTitle),
		fwCheckWaitDurationBlock(),
		widget.NewSeparator(),
		fwTabs(),
		widget.NewSeparator(),
	)
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
		strSlice, err := strList.Get()
		if err != nil {
			fyne.LogError("addBtn.strList.Get", err)
		}
		if text != "" && !slices.Contains(strSlice, text) {
			if err := strList.Append(text); err != nil {
				fyne.LogError("stringListBlock.strList.Append", err)
			}
		}
		entry.SetText("")
	})
	rmBtn := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		text := entry.Text
		strSlice, err := strList.Get()
		if err != nil {
			fyne.LogError("rmBtn.strList.Get", err)
		}
		if text != "" && slices.Contains(strSlice, text) {
			if err := strList.Remove(text); err != nil {
				fyne.LogError("stringListBlock.strList.Remove", err)
			}
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
		selectedItemText, err := strList.GetValue(id)
		if err != nil {
			fyne.LogError("list.OnSelected.strList.GetValue", err)
		}
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
	titleLabel := widget.NewLabel(title)
	titleLabel.Wrapping = fyne.TextWrapWord
	return container.NewBorder(nil, nil, resetBtn(reset), nil, titleLabel)
}

func durationBlock(
	title string, min, max time.Duration, data binding.String, validator func(v string) error, setDefault func(),
) *fyne.Container {
	e := NewStrPrefsEntry(data)
	e.Validator = validator
	return container.NewVBox(
		prefTitleWithResetBtn(
			fmt.Sprintf("%s (от %s до %s включительно)", title, min, max),
			setDefault,
		), e,
	)
}
