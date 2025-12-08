package ui

import (
	"pc_security_test/command"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type searchBlockParams struct {
	cols        []searchBlockTableCol
	commandFunc func() command.Request
	awaitFunc   func() command.SearchResponse
}

type searchBlockTableCol struct {
	title string
	rows  binding.StringList
}

type searchBlock struct {
	searchBtn   *widget.Button
	commandFunc func() command.Request
	clearBtn    *widget.Button
	cols        []searchBlockTableCol
	tbl         *widget.Table
	templateLen int
	awaitFunc   func() command.SearchResponse
	results     *strListWithImportance
}

func newSearchBlock(params searchBlockParams) *searchBlock {
	block := &searchBlock{}

	block.searchBtn = widget.NewButtonWithIcon("Поиск", theme.SearchIcon(), block.search)
	block.clearBtn = widget.NewButtonWithIcon("Сбросить", theme.ContentClearIcon(), block.clearResults)

	block.commandFunc = params.commandFunc
	block.awaitFunc = params.awaitFunc

	block.results = newStrListWithImportance()

	block.cols = params.cols
	block.recountTemplateLen()
	block.tbl = widget.NewTable(block.tableLength(params.cols), block.createItem, block.updateItem)
	block.tbl.CreateHeader = block.createItem
	block.tbl.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		template.(*widget.Label).SetText(block.cols[id.Col].title)
	}
	block.tbl.ShowHeaderRow = true
	block.addTemplateLenListeners()

	return block
}

func (b *searchBlock) tableLength(tcs []searchBlockTableCol) func() (rows int, cols int) {
	return func() (rows int, cols int) {
		cols = len(tcs)
		for _, tc := range tcs {
			if tc.rows.Length() > rows {
				rows = tc.rows.Length()
			}
		}
		return
	}
}

func (b *searchBlock) addTemplateLenListeners() {
	for _, tc := range b.cols {
		tc.rows.AddListener(binding.NewDataListener(func() {
			b.recountTemplateLen()
			b.tbl.Refresh()
		}))
	}
}

func (b *searchBlock) recountTemplateLen() {
	b.templateLen = 0
	for _, tc := range b.cols {
		if utf8.RuneCountInString(tc.title) > b.templateLen {
			b.templateLen = utf8.RuneCountInString(tc.title)
		}
		for i := range tc.rows.Length() {
			v, _ := tc.rows.GetValue(i)
			if utf8.RuneCountInString(v) > b.templateLen {
				b.templateLen = utf8.RuneCountInString(v)
			}
		}
	}
}

func (b *searchBlock) makeTemplate() string {
	var sb strings.Builder
	for range b.templateLen {
		sb.WriteRune('o')
	}
	return sb.String()
}

func (b *searchBlock) createItem() fyne.CanvasObject {
	return widget.NewLabel(b.makeTemplate())
}

func (b *searchBlock) updateItem(tci widget.TableCellID, co fyne.CanvasObject) {
	col := b.cols[tci.Col]
	rowText, _ := col.rows.GetValue(tci.Row)
	co.(*widget.Label).Text = rowText
	co.(*widget.Label).Importance = b.results.Get(rowText)
	co.(*widget.Label).Refresh()
}

func (b *searchBlock) clearResults() {
	b.results.Clear()
	b.tbl.Refresh()
}

func (b *searchBlock) search() {
	b.clearResults()
	go command.AddToQueue(b.commandFunc())
	go b.awaitAndUpdateUI()
}

func (b *searchBlock) awaitAndUpdateUI() {
	res := b.awaitFunc()

	fyne.Do(func() {
		b.fillResults(res)
		b.tbl.Refresh()
	})
}

func (b *searchBlock) fillResults(res command.SearchResponse) {
	for _, col := range b.cols {
		strs, _ := col.rows.Get()
		for _, s := range strs {
			if _, ok := res.Found[s]; ok {
				b.results.Set(s, widget.SuccessImportance)
			} else {
				b.results.Set(s, widget.DangerImportance)
			}
		}
	}
}

func (b *searchBlock) getContainer() *fyne.Container {
	tblScroll := container.NewScroll(b.tbl)
	tblScroll.SetMinSize(fyne.NewSize(0, 210))
	return container.NewVBox(container.NewHBox(b.searchBtn, b.clearBtn), tblScroll)
}
