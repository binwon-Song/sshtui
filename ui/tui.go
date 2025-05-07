package ui

import (
	"log"
	"sshtui/config"
	"sshtui/ui/components"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartUI(cfg *config.Config) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// 서버 리스트 생성
	serverList := components.NewServerList(app, cfg)

	// 하단 네비게이션 바
	navbar := components.NewNavbar()

	// 메인 레이아웃
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(serverList.GetPrimitive(), 0, 1, true).
		AddItem(navbar, 3, 0, false)

	pages.AddPage("main", mainLayout, true, true)

	// 모달 표시 함수들
	showAddServerModal := func() {
		serverForm := components.NewServerForm(app, pages, func(newServer config.Server) {
			cfg.Servers = append(cfg.Servers, newServer)
			if err := config.SaveConfig("", cfg); err != nil {
				log.Printf("설정 저장 실패: %v", err)
			}
			serverList.AddServer(newServer)
			pages.RemovePage("modal")
			app.SetFocus(serverList.GetPrimitive())
		})

		modalFlex := components.CreateModalFlex(serverForm.GetPrimitive())
		pages.AddPage("modal", modalFlex, false, true)
		app.SetFocus(serverForm.GetPrimitive())
	}

	showHelpModal := func() {
		helpModal := components.NewHelpModal(func() {
			pages.RemovePage("help")
			app.SetFocus(serverList.GetPrimitive())
		})
		pages.AddPage("help", helpModal, false, true)
		app.SetFocus(helpModal)
	}

	// 키 바인딩 설정
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentPage, _ := pages.GetFrontPage()
		if currentPage == "modal" || currentPage == "help" {
			if event.Key() == tcell.KeyEscape {
				pages.RemovePage(currentPage)
				app.SetFocus(serverList.GetPrimitive())
				return nil
			}
			switch event.Rune() {
			case 'a', 'q', 'h':
				return nil // 모달이 열려있을 때는 a, q, h 키 무시
			}
		} else {
			switch event.Rune() {
			case 'a':
				showAddServerModal()
				return nil
			case 'h':
				showHelpModal()
				return nil
			case 'q':
				app.Stop()
				return nil
			}
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
