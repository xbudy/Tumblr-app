package main

import "fyne.io/fyne/v2/widget"

type Logg struct {
	Entry *widget.Entry
}

func (log Logg) WrtiteLog(text string) {
	log.Entry.SetText(log.Entry.Text + "\n" + text)
	log.Entry.Refresh()
}
