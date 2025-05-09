package components

import (
	"fmt"
	"log"
	"sshtui/config"
	"sshtui/ssh"

	"github.com/rivo/tview"
)

type ServerList struct {
	list *tview.List
	app  *tview.Application
	cfg  *config.Config
}

func NewServerList(app *tview.Application, cfg *config.Config) *ServerList {
	sl := &ServerList{
		list: tview.NewList(),
		app:  app,
		cfg:  cfg,
	}
	sl.initialize()
	return sl
}

func (sl *ServerList) initialize() {
	sl.list.SetBorder(true).SetTitle("SSH 서버 목록")
	sl.refreshItems()
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

func (sl *ServerList) GetPrimitive() *tview.List {
	return sl.list
}
