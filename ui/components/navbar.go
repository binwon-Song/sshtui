package components

import "github.com/rivo/tview"

func NewNavbar() *tview.TextView {
	return tview.NewTextView().
		SetText("[yellow]서버 추가[white] [a[]  [yellow]도움말[white] [h[]  [yellow]종료[white] [q[]  [yellow]나가기[white] [ESC[]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
}
