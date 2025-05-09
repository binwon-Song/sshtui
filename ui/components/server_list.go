package components

import (
	"fmt"
	"log"
	"sshtui/config"
	"sshtui/ssh"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ServerList struct {
	list  *tview.List
	app   *tview.Application
	cfg   *config.Config
	pages *tview.Pages
}

func NewServerList(app *tview.Application, cfg *config.Config, pages *tview.Pages) *ServerList {
	sl := &ServerList{
		list:  tview.NewList(),
		app:   app,
		cfg:   cfg,
		pages: pages,
	}
	sl.initialize()
	return sl
}

func (sl *ServerList) initialize() {
	sl.list.SetBorder(true).SetTitle("SSH 서버 목록")
	sl.refreshItems()

	confirmModal := NewConfirmModal(sl.app, sl.pages)

	// 키 입력 핸들러 추가
	sl.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'd':
			currentItem := sl.list.GetCurrentItem()
			if currentItem >= 0 && currentItem < len(sl.cfg.Servers) {
				confirmModal.Show(sl, currentItem)
				return nil
			}
		}
		return event
	})
}

func (sl *ServerList) refreshItems() {
	sl.list.Clear()
	for _, server := range sl.cfg.Servers {
		serverCopy := server
		label := fmt.Sprintf("%s (%s@%s:%d)", server.Name, server.User, server.Host, server.Port)
		sl.list.AddItem(label, "", 0, func() {
			sl.app.Suspend(func() {
				if err := ssh.Connect(serverCopy); err != nil {
					log.Printf("SSH 접속 실패: %v", err)
				}
			})
		})
	}
}

func (sl *ServerList) AddServer(server config.Server) {
	sl.cfg.Servers = append(sl.cfg.Servers, server)

	// YAML 파일에 저장
	if err := config.SaveConfig(sl.cfg); err != nil {
		log.Printf("설정 저장 실패: %v", err)
	}

	sl.refreshItems()
}

func (sl *ServerList) DeleteServer(index int) {
	if index < 0 || index >= len(sl.cfg.Servers) {
		return
	}

	serverToDelete := sl.cfg.Servers[index]
	sl.cfg.Servers = append(sl.cfg.Servers[:index], sl.cfg.Servers[index+1:]...)

	// YAML 파일에 저장
	if err := config.SaveConfig(sl.cfg); err != nil {
		log.Printf("서버 '%s' 삭제 실패: %v", serverToDelete.Name, err)
		return
	}

	log.Printf("서버 '%s' 삭제됨", serverToDelete.Name)
	sl.refreshItems()
}

func (sl *ServerList) GetPrimitive() *tview.List {
	return sl.list
}
