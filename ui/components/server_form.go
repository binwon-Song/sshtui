package components

import (
	"fmt"
	"sshtui/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ServerForm struct {
	form   *tview.Form
	app    *tview.Application
	pages  *tview.Pages
	onSave func(config.Server)
}

func NewServerForm(app *tview.Application, pages *tview.Pages, onSave func(config.Server)) *ServerForm {
	sf := &ServerForm{
		form:   tview.NewForm(),
		app:    app,
		pages:  pages,
		onSave: onSave,
	}
	sf.initialize()
	return sf
}

func (sf *ServerForm) initialize() {
	var name, user, host string
	var port int

	sf.form.AddInputField("Name", "", 20, nil, func(text string) { name = text })
	sf.form.AddInputField("User", "", 20, nil, func(text string) { user = text })
	sf.form.AddInputField("Host", "", 20, nil, func(text string) { host = text })
	sf.form.AddInputField("Port", "22", 20, nil, func(text string) {
		if text == "" { // 사용자가 아무것도 입력하지 않은 경우
			port = 22
		} else {
			fmt.Sscanf(text, "%d", &port)
		}
	})

	sf.form.AddButton("추가", func() {
		newServer := config.Server{
			Name: name,
			User: user,
			Host: host,
			Port: port,
		}
		sf.onSave(newServer)
		config.SaveConfig("", nil) // 서버 추가 후 설정 저장
	})

	sf.form.AddButton("취소", func() {
		sf.pages.RemovePage("modal")
	})

	sf.form.SetBorder(true).
		SetTitle("서버 추가").
		SetTitleAlign(tview.AlignLeft)

	sf.form.SetInputCapture(sf.handleKeyEvents)
}

func (sf *ServerForm) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape: // ESC 키를 눌렀을 때
		sf.pages.RemovePage("modal")
		return nil
	case tcell.KeyUp: // 위 방향키를 눌렀을 때
		for i := 1; i < sf.form.GetFormItemCount(); i++ {
			if sf.form.GetFormItem(i).HasFocus() {
				sf.app.SetFocus(sf.form.GetFormItem(i - 1))
				return nil
			}
		}
	case tcell.KeyDown: // 아래 방향키를 눌렀을 때
		for i := 0; i < sf.form.GetFormItemCount()-1; i++ {
			if sf.form.GetFormItem(i).HasFocus() {
				sf.app.SetFocus(sf.form.GetFormItem(i + 1))
				return nil
			}
		}
	}
	return event
}

func (sf *ServerForm) GetPrimitive() *tview.Form {
	return sf.form
}

func CreateModalFlex(form tview.Primitive) *tview.Flex {
	modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	modalWidth := 40
	modalHeight := 12
	leftPadding := 5
	topPadding := 2

	return modalFlex.
		AddItem(nil, topPadding, 0, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, leftPadding, 0, false).
				AddItem(form, modalWidth, 0, true).
				AddItem(nil, leftPadding, 0, false),
			modalHeight, 0, true,
		).
		AddItem(nil, topPadding, 0, false)
}
