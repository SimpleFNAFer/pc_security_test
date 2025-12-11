package ui

import (
	"fmt"
	"pc_security_test/command"
	"pc_security_test/preferences"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const (
	fwSearchTitle = "## Проверка наличия межсетевого экрана"
	fwSearchDesc  = `Данный раздел предназначен для поиска исполняемых файлов (в PATH), а также файловых путей.

Список доступных для поиска файлов и путей можно отредактировать в настройках в разделе "Межсетевой экран".

Для начала поиска нажмите "Поиск". Найденные и не найденные объекты подсветятся соответствующими цветами.
Для сброса результатов нажмите "Сбросить"
	`
)

type fwSearchBlock struct {
	sb *searchBlock
}

func newSearchFWBlock() *fwSearchBlock {
	block := &fwSearchBlock{}
	block.sb = newSearchBlock(searchBlockParams{
		cols: []searchBlockTableCol{
			{
				title: binariesTitle,
				rows:  preferences.FWBinaries,
			},
			{
				title: filepathsTitle,
				rows:  preferences.FWFilePaths,
			},
		},
		commandFunc: func() command.Request {
			return command.FindFWRequest{
				ID: uuid.New(),
			}
		},
		awaitFunc: command.AwaitFindFWResponse,
	})

	return block
}

func (b *fwSearchBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(fwSearchDesc)
	desc.Wrapping = fyne.TextWrapWord
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(fwSearchTitle), desc,
		b.sb.getContainer(),
	)
}

const (
	checkFWTitle = "## Проверка работы межсетевого экрана"
	//nolint:lll
	checkFWDesc = `Данный раздел предназначен для проверки работы межсетевого экрана путём обращения к ресурсу по определённому протоколу
(опционально - по определённому порту).

Ресурс должен быть не пустым и являться валидным доменным именем, IPv4 или IPv6 адресом.
Доступно два протокола: tcp и udp.
Порт является необязательным параметром. При выборе протокола tcp порт по умолчанию - 443. В случае udp - 53.
Длительность ожидания подключения можно изменить в настройках в разделе "Межсетевой экран"

Для начала проверки нажмите "Проверить". Результат отобразится под полем ввода, а также будет сохранён в истории.
	`
)

type checkFWBlock struct {
	protocol     *widget.Select
	hostInput    *widget.Entry
	portInput    *widget.Entry
	resultOutput *widget.Label
}

func newCheckFWBlock() *checkFWBlock {
	block := &checkFWBlock{}

	protocol := widget.NewSelect(preferences.Protocols, nil)
	protocol.SetSelectedIndex(0)

	hostInput := widget.NewEntry()
	preferences.PingDefaultHost.AddListener(binding.NewDataListener(func() {
		ph, err := preferences.PingDefaultHost.Get()
		if err != nil {
			fyne.LogError("newCheckFWBlock.PingDefaultHost.Get", err)
		}
		hostInput.SetPlaceHolder(ph)
	}))
	hostInput.Validator = func(s string) error {
		if s == "" {
			s = hostInput.PlaceHolder
		}
		return preferences.HostValidator(s)
	}

	portInput := widget.NewEntry()
	portInput.Validator = func(s string) error {
		if portInput.Text == "" {
			return nil
		}
		return preferences.MinMaxIntValidator(0, 65535)(s)
	}

	resultOutput := widget.NewLabel("")
	resultOutput.Wrapping = fyne.TextWrapWord
	resultOutput.Hide()

	block.protocol = protocol
	block.hostInput = hostInput
	block.portInput = portInput
	block.resultOutput = resultOutput

	return block
}

func (b *checkFWBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(checkFWDesc)
	desc.Wrapping = fyne.TextWrapWord

	form := widget.NewForm(
		widget.NewFormItem("Протокол", b.protocol),
		widget.NewFormItem("Хост", b.hostInput),
		widget.NewFormItem("*Порт", b.portInput),
	)
	form.SubmitText = "Проверить"
	form.OnSubmit = b.onCheckButtonClick

	return container.NewVBox(
		widget.NewRichTextFromMarkdown(checkFWTitle), desc,
		form,
		b.resultOutput,
	)
}

func (b *checkFWBlock) onCheckButtonClick() {
	host := b.hostInput.Text
	if host == "" {
		host = b.hostInput.PlaceHolder
	}

	go command.AddToQueue(command.TestFWRequest{
		ID:       uuid.New(),
		Protocol: preferences.Protocol(b.protocol.Selected),
		Host:     host,
		Port:     b.portInput.Text,
	})
	go b.awaitAndUpdateUI()
}

func (b *checkFWBlock) awaitAndUpdateUI() {
	res := command.AwaitTestFWResponse()

	var resStr strings.Builder
	resStr.WriteString("к хосту " + res.Host)
	if res.Port != "" {
		resStr.WriteString(fmt.Sprintf(" (порт: %s)", res.Port))
	}
	resStr.WriteString(fmt.Sprintf(" по протоколу %s", res.Protocol))

	fyne.Do(func() {
		switch {
		case res.Error != nil:
			b.resultOutput.Text = res.Error.Error()
			b.resultOutput.Importance = widget.DangerImportance
		case res.Unavailable:
			b.resultOutput.Text = fmt.Sprintf("МЭ работает, нет доступа %s", resStr.String())
			b.resultOutput.Importance = widget.SuccessImportance
		default:
			b.resultOutput.Text = fmt.Sprintf("МЭ не работает, есть доступ %s", resStr.String())
			b.resultOutput.Importance = widget.WarningImportance
		}
		b.resultOutput.Refresh()
		b.resultOutput.Show()
	})
}
