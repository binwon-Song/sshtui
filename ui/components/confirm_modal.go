package components

import (
	"github.com/rivo/tview"
)

type ConfirmModal struct {
	modal *tview.Modal
	app   *tview.Application
	pages *tview.Pages
}

func NewConfirmModal(app *tview.Application, pages *tview.Pages) *ConfirmModal {
	return &ConfirmModal{
		modal: tview.NewModal(),
		app:   app,
		pages: pages,
	}
}

func (cm *ConfirmModal) Show(serverList *ServerList, index int) {
	serverToDelete := serverList.cfg.Servers[index]
	cm.modal.ClearButtons() // 기존 버튼들을 초기화
	cm.modal.
		SetText(serverToDelete.Name + " 서버를 삭제하시겠습니까?").
		AddButtons([]string{"삭제", "취소"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "삭제" {
				serverList.DeleteServer(index)
			}
			cm.app.SetRoot(cm.pages, true)
		})

	cm.app.SetRoot(cm.modal, true)
}
