package main

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func CreatePostsListWindow(mdb MongoDb, a fyne.App, blog string) fyne.Window {
	w := a.NewWindow(blog + "'s posts")
	w.Resize(fyne.NewSize(780, 400))
	data := BindPosts(mdb, blog)

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
				pid, _ := t.Get()
				u, _ := url.Parse(fmt.Sprintf("https://%v.tumblr.com/post/%v", blog, pid))
				a.OpenURL(u)
			}

		})
	label1 := widget.NewLabel("Id")
	value1 := widget.NewLabel("..")
	label2 := widget.NewLabel("Images")
	value2 := widget.NewLabel("...")
	label3 := widget.NewLabel("Url")
	value3 := widget.NewLabel("...")
	label4 := widget.NewLabel("Date")
	value4 := widget.NewLabel("...")
	list.OnSelected = func(id widget.ListItemID) {
		pid, _ := data.GetValue(id)
		post := mdb.getPost(blog, pid)
		value1.Text = post.Id
		value1.TextStyle = fyne.TextStyle{Bold: true}
		value2.Text = fmt.Sprint(len(post.Medias))
		value3.Text = fmt.Sprintf("https://%v.tumblr.com/post/%v", blog, pid)
		value4.Text = fmt.Sprint(post.Timestamp)
		// label2.TextStyle = fyne.TextStyle{Bold: true}
		value1.Refresh()
		value2.Refresh()
		value3.Refresh()
		value4.Refresh()
	}
	grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2, label3, value3, label4, value4)
	//c := container.NewVBox(label1, label2)
	w.SetContent(container.NewHSplit(
		list,
		grid,
	))
	return w
}
