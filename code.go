package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func f() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Data")

	data := binding.BindStringList(
		&[]string{"Item 1", "Item 2", "Item 3"},
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {

			return container.NewBorder(nil, nil, nil, widget.NewButton("Open", nil),
				widget.NewLabel("item x.y"))
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			t := i.(binding.String)
			f := o.(*fyne.Container).Objects[0]
			f.(*widget.Label).Bind(i.(binding.String))
			btn := o.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = func() {
				log.Println(t.Get())
			}

		})

	add := widget.NewButton("Append", func() {
		val := fmt.Sprintf("Item %d", data.Length()+1)
		data.Append(val)
	})
	myWindow.SetContent(container.NewBorder(nil, add, nil, nil, list))
	myWindow.ShowAndRun()
}
