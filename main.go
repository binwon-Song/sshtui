package main

import (
	"log"
	"sshtui/config"
	"sshtui/ui"
)

func main() {
	cfg, err := config.LoadConfig("")

	if err != nil {
		log.Fatalf("서버 설정 로딩 실패: %v", err)
	}

	ui.StartUI(cfg)
}
