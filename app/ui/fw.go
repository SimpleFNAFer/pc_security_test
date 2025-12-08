package ui

import (
	"pc_security_test/command"
	"pc_security_test/preferences"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

const (
	fwSearchTitle = "## Проверка наличия межсетевого экрана"
	fwSearchDesc  = `
Данный раздел предназначен для поиска исполняемых файлов (в PATH), а также файловых путей.

Список доступных для поиска файлов и путей можно отредактировать в настройках в разделе "Межсетевой экран".

Для начала поиска нажмите "Поиск". Найденные и не найденные объекты подсветятся соответствующими цветами.
Для сброса результатов нажмите "Сбросить"
	`
)

type findFWBlock struct {
	sb *searchBlock
}

func newSearchFWBlock() *findFWBlock {
	block := &findFWBlock{}
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

func (b *findFWBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(fwSearchDesc)
	desc.Wrapping = fyne.TextWrapWord
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(fwSearchTitle),
		desc,
		b.sb.getContainer(),
	)
}
