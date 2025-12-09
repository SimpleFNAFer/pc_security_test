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
	avSearchTitle = "## Проверка наличия антивируса"
	avSearchDesc  = `Данный раздел предназначен для поиска исполняемых файлов (в PATH), а также файловых путей.

Список доступных для поиска файлов и путей можно отредактировать в настройках в разделе "Антивирус".

Для начала поиска нажмите "Поиск". Найденные и не найденные объекты подсветятся соответствующими цветами.
Для сброса результатов нажмите "Сбросить"
	`
)

type avSearchBlock struct {
	sb *searchBlock
}

func newSearchAVBlock() *avSearchBlock {
	block := &avSearchBlock{}
	block.sb = newSearchBlock(searchBlockParams{
		cols: []searchBlockTableCol{
			{
				title: binariesTitle,
				rows:  preferences.AVBinaries,
			},
			{
				title: filepathsTitle,
				rows:  preferences.AVFilePaths,
			},
		},
		commandFunc: func() command.Request {
			return command.FindAVRequest{
				ID: uuid.New(),
			}
		},
		awaitFunc: command.AwaitFindAVResponse,
	})

	return block
}

func (b *avSearchBlock) getContainer() *fyne.Container {
	desc := widget.NewLabel(avSearchDesc)
	desc.Wrapping = fyne.TextWrapWord
	return container.NewVBox(
		widget.NewRichTextFromMarkdown(avSearchTitle), desc,
		b.sb.getContainer(),
	)
}
