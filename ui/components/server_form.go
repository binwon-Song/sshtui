package components

import (
	"fmt"
	"sshtui/config"
	"sshtui/crypto"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ServerForm struct {
	form   *tview.Form
	app    *tview.Application
	pages  *tview.Pages
	onSave func(config.Server)
	debug  func(string)   // 디버그 메시지를 표시하기 위한 함수 추가
	cfg    *config.Config // Config 객체 추가
}

/*
 * NewServerForm creates a new server form.
 * @param app *tview.Application
 * @param pages *tview.Pages
 * @param onSave function callback to save the server configuration parameter is Server
 * @param cfg *config.Config
 * @return ServerForm
 */
func NewServerForm(app *tview.Application, pages *tview.Pages, onSave func(config.Server), cfg *config.Config) *ServerForm {
	sf := &ServerForm{
		form:   tview.NewForm(),
		app:    app,
		pages:  pages,
		onSave: onSave,
		cfg:    cfg,
		debug: func(msg string) { // 디버그 메시지를 표시하는 함수
			modal := tview.NewModal().
				SetText(msg).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					pages.RemovePage("debug")
				})
			pages.AddPage("debug", modal, true, true)
		},
	}
	sf.initialize()
	return sf
}

func (sf *ServerForm) initialize() {
	var name, user, host, password string
	var port int = 22

	sf.form.AddInputField("Name", "", 20, nil, func(text string) { name = text })
	sf.form.AddInputField("User", "", 20, nil, func(text string) { user = text })
	sf.form.AddInputField("Host", "", 20, nil, func(text string) { host = text })
	sf.form.AddInputField("Port", "22", 20, nil, func(text string) {
		if text != "" {
			fmt.Sscanf(text, "%d", &port)
		}
	})
	sf.form.AddPasswordField("Password (선택)", "", 20, '*', func(text string) { password = text })

	sf.form.AddButton("추가", func() {
		var encryptedPass string
		if password != "" {
			// 비밀번호가 입력된 경우에만 암호화
			encrypted, err := crypto.Encrypt(password)
			if err != nil {
				sf.debug(fmt.Sprintf("비밀번호 암호화 실패: %v", err))
				return
			}
			encryptedPass = encrypted
		}

		newServer := config.Server{
			Name:          name,
			User:          user,
			Host:          host,
			Port:          port,
			EncryptedPass: encryptedPass,
		}

		sf.onSave(newServer)
		sf.pages.RemovePage("modal")
	})

	sf.form.AddButton("취소", func() {
		sf.pages.RemovePage("modal")
	})

	sf.form.
		SetBorder(true).
		SetTitle("Add Server").
		SetTitleAlign(tview.AlignCenter)

	sf.form.SetButtonsAlign(tview.AlignCenter)
	sf.form.SetInputCapture(sf.handleKeyEvents)
}

func (sf *ServerForm) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		sf.pages.RemovePage("modal")
		return nil
	case tcell.KeyUp:
		// 현재 포커스된 아이템의 인덱스 찾기
		for i := 0; i < sf.form.GetFormItemCount(); i++ {
			if sf.form.GetFormItem(i).HasFocus() {
				if i > 0 { // 첫 번째 아이템이 아닌 경우
					sf.app.SetFocus(sf.form.GetFormItem(i - 1))
				}
				return nil
			}
		}
	case tcell.KeyDown:
		// 현재 포커스된 아이템의 인덱스 찾기
		for i := 0; i < sf.form.GetFormItemCount(); i++ {
			if sf.form.GetFormItem(i).HasFocus() {
				if i < sf.form.GetFormItemCount()-1 { // 마지막 아이템이 아닌 경우
					sf.app.SetFocus(sf.form.GetFormItem(i + 1))
				}
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
	modalHeight := 15
	rightPadding := 2
	topPadding := 2

	return modalFlex.
		AddItem(nil, topPadding, 0, false).
		AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(form, modalWidth, 0, true).
				AddItem(nil, rightPadding, 0, false),
			modalHeight, 0, true,
		).
		AddItem(nil, topPadding, 0, false)
}
