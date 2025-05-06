package ui

import (
	"fmt"
	"log"
	"sshtui/config"
	"sshtui/ssh"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func StartUI(cfg *config.Config) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// 서버 리스트 생성
	serverList := tview.NewList()
	serverList.
		SetBorder(true).
		SetTitle("SSH 서버 목록")

	// 서버 목록 초기화
	for _, server := range cfg.Servers {
		serverCopy := server
		label := fmt.Sprintf("%s (%s@%s:%d)", server.Name, server.User, server.Host, server.Port)
		serverList.AddItem(label, "", 0, func() {
			app.Suspend(func() {
				if err := ssh.Connect(serverCopy); err != nil {
					log.Printf("SSH 접속 실패: %v", err)
				}
			})
		})
	}

	// 하단 네비게이션 바
	navbar := tview.NewTextView().
		SetText("[yellow]서버 추가[white] [a[]  [yellow]도움말[white] [h[]  [yellow]종료[white] [q[]  [yellow]나가기[white] [ESC[]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// 메인 레이아웃
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(serverList, 0, 1, true).
		AddItem(navbar, 3, 0, false)

	pages.AddPage("main", mainLayout, true, true)

	// 모달 생성 함수
	showAddServerModal := func() {
		modalFlex := tview.NewFlex().
			SetDirection(tview.FlexRow)

		form := createServerAddForm(app, pages, serverList, cfg, func(newServer config.Server) {
			// 리스트에 추가
			cfg.Servers = append(cfg.Servers, newServer)
			label := fmt.Sprintf("%s (%s@%s:%d)", newServer.Name, newServer.User, newServer.Host, newServer.Port)
			serverList.AddItem(label, "", 0, func() {
				app.Suspend(func() {
					if err := ssh.Connect(newServer); err != nil {
						log.Printf("SSH 접속 실패: %v", err)
					}
				})
			})
			pages.RemovePage("modal")
			app.SetFocus(serverList)
		})

		// 모달 크기 및 위치 조정
		modalWidth := 40
		modalHeight := 12

		// 고정된 여백 사용
		leftPadding := 5
		topPadding := 2

		// 모달 레이아웃 구성
		modalFlex.
			AddItem(nil, topPadding, 0, false).
			AddItem(
				tview.NewFlex().
					AddItem(nil, leftPadding, 0, false).
					AddItem(form, modalWidth, 0, true).
					AddItem(nil, leftPadding, 0, false),
				modalHeight, 0, true,
			).
			AddItem(nil, topPadding, 0, false)

		pages.AddPage("modal", modalFlex, false, true)
		app.SetFocus(form)
	}

	// 도움말 모달 표시 함수
	showHelpModal := func() {
		helpText := "SSH 서버 관리 프로그램\n\n" +
			"[yellow]키보드 단축키:[white]\n" +
			"[a[] : 새로운 서버 추가\n" +
			"[h[] : 이 도움말 표시\n" +
			"[q[] : 프로그램 종료\n" +
			"[ESC[] : 현재 화면 닫기\n\n" +
			"[yellow]입력 폼에서:[white]\n" +
			"[↑][↓] : 필드 간 이동\n" +
			"[Tab[] : 다음 필드로 이동"

		modal := tview.NewModal().
			SetText(helpText).
			AddButtons([]string{"확인"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("help")
				app.SetFocus(serverList)
			})

		pages.AddPage("help", modal, false, true)
		app.SetFocus(modal)
	}

	// 키 바인딩 설정
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentPage, _ := pages.GetFrontPage()
		if currentPage == "modal" || currentPage == "help" {
			if event.Key() == tcell.KeyEscape {
				pages.RemovePage(currentPage)
				app.SetFocus(serverList)
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

func createServerAddForm(app *tview.Application, pages *tview.Pages, serverList *tview.List, cfg *config.Config, onSave func(config.Server)) tview.Primitive {
	form := tview.NewForm()

	var name, user, host string
	var port int

	form.AddInputField("Name", "", 20, nil, func(text string) { name = text })
	form.AddInputField("User", "", 20, nil, func(text string) { user = text })
	form.AddInputField("Host", "", 20, nil, func(text string) { host = text })
	form.AddInputField("Port", "22", 20, nil, func(text string) {
		fmt.Sscanf(text, "%d", &port)
	})

	form.AddButton("추가", func() {
		newServer := config.Server{
			Name: name,
			User: user,
			Host: host,
			Port: port,
		}
		onSave(newServer)
	})

	form.AddButton("취소", func() {
		pages.SwitchToPage("main")
		app.SetFocus(serverList)
	})

	form.SetBorder(true).
		SetTitle("서버 추가").
		SetTitleAlign(tview.AlignLeft)

	// 폼에 키보드 이벤트 핸들러 추가
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			pages.SwitchToPage("main")
			app.SetFocus(serverList)
			return nil
		case tcell.KeyUp:
			// 현재 포커스된 아이템의 인덱스를 찾기
			for i := 1; i < form.GetFormItemCount(); i++ {
				if form.GetFormItem(i).HasFocus() {
					app.SetFocus(form.GetFormItem(i - 1))
					return nil
				}
			}
		case tcell.KeyDown:
			// 현재 포커스된 아이템의 인덱스를 찾기
			for i := 0; i < form.GetFormItemCount()-1; i++ {
				if form.GetFormItem(i).HasFocus() {
					app.SetFocus(form.GetFormItem(i + 1))
					return nil
				}
			}
		}
		return event
	})

	return form
}
