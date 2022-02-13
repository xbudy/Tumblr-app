package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	mainWindow := a.NewWindow("Tumblr")
	//InitDb
	cl, ctx := InitMongoDb()
	Mongodb := MongoDb{Client: cl, Ctx: ctx}
	defer Mongodb.Client.Disconnect(Mongodb.Ctx)
	//Blogname
	BlogLabel := widget.NewLabel("Blog")
	BlogEntry := widget.NewEntry()
	BlogEntry.SetPlaceHolder("Enter Blog Name")
	//Logging Space
	loggingTitle := widget.NewLabel("LOGS")
	loggingTitle.Alignment = fyne.TextAlignCenter
	logginSpace := widget.NewMultiLineEntry()
	logSpace := container.NewVBox(loggingTitle, layout.NewSpacer(), logginSpace)
	//Init Logs
	logg := Logg{Entry: logginSpace}
	//Posts Count
	LabelPostsCount := widget.NewLabel("Posts Scraped")
	PostsCount := widget.NewLabel("..")
	PostsCountForm := container.New(layout.NewFormLayout(), LabelPostsCount, PostsCount)
	//ShowPostsButton
	ShowPostsBtn := widget.NewButton("Show Posts", func() {
		if ThereisBlogName(BlogEntry) {
			w := CreatePostsListWindow(Mongodb, a, BlogEntry.Text)
			w.Show()
		} else {
			dialog.ShowError(fmt.Errorf("please enter blog name"), mainWindow)
		}

	})
	//Extracting Posts Button
	ExtractPostsBtn := widget.NewButton("Extract Posts", func() {
		ExtractWindow := CreateExtractPostsWindow(Mongodb, a, BlogEntry.Text)
		ExtractWindow.Show()
	})
	//Horizontal Layout For Buttons Show and extract
	HboxButtons := container.NewHBox(ShowPostsBtn, layout.NewSpacer(), ExtractPostsBtn)
	//Blog Load
	Form := container.New(layout.NewFormLayout(), BlogLabel, BlogEntry) //Form Blog Name

	LoadBlog := widget.NewButton("Load Blog", func() {
		if !ThereisBlogName(BlogEntry) {
			dialog.ShowError(fmt.Errorf("please enter blog name"), mainWindow)
		} else {
			Count, _ := Mongodb.GetPostsCount(BlogEntry.Text)
			PostsCount.SetText(fmt.Sprint(Count))
			PostsCount.Refresh()
			logg.WrtiteLog("Loadding " + BlogEntry.Text)
		}
	})
	//
	Vbox := container.NewVBox(Form, LoadBlog, PostsCountForm, HboxButtons)
	//
	MainVbox := container.NewVSplit(Vbox, logSpace)
	mainWindow.SetContent(MainVbox)
	mainWindow.SetMaster()
	mainWindow.ShowAndRun()
}

func BindPosts(mdb MongoDb, blog string) binding.ExternalStringList {
	dataP := mdb.GetData(blog)
	var ids []string
	for _, p := range dataP {
		ids = append(ids, p.Id)
	}
	data := binding.BindStringList(
		&ids,
	)
	return data
}
func ThereisBlogName(entry *widget.Entry) bool {

	return len(entry.Text) != 0
}
