package components

import "github.com/rivo/tview"

func NewNavbar() *tview.TextView {
	return tview.NewTextView().
		SetText("[yellow]Add[white] [a[] [yellow]Delete[white] [d[] [yellow]HELP[white] [h[]  [yellow]Quit[white] [q[]  [yellow]CLOSE[white] [ESC[]").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
}
