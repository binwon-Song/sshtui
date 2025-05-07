package components

import "github.com/rivo/tview"

func NewHelpModal(onClose func()) *tview.Modal {
	helpText := "SSH 서버 관리 프로그램\n\n" +
		"[yellow]키보드 단축키:[white]\n" +
		"[a[] : 새로운 서버 추가\n" +
		"[h[] : 이 도움말 표시\n" +
		"[q[] : 프로그램 종료\n" +
		"[ESC[] : 현재 화면 닫기\n\n" +
		"[yellow]입력 폼에서:[white]\n" +
		"[↑][↓] : 필드 간 이동\n" +
		"[Tab[] : 다음 필드로 이동"

	return tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"확인"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			onClose()
		})
}
