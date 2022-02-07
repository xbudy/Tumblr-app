package main

import (
	"fmt"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
func CreateExtractPostsWindow(mdb MongoDb, a fyne.App, blog string) fyne.Window {
	w := a.NewWindow("Extract " + blog + " Posts")
	w.Resize(fyne.NewSize(780, 400))
	infoLabel := widget.NewLabel("Info")
	InfoTitle := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), infoLabel, layout.NewSpacer())
	LastPageLabel := widget.NewLabel("Last Page")
	LastPageValue := widget.NewLabel("Test Value")
	PostDownloadedLabel := widget.NewLabel("Downloaded Posts")
	PostDownloadedValue := widget.NewLabel("300")
	InfoForm := container.New(layout.NewFormLayout(), LastPageLabel, LastPageValue, PostDownloadedLabel, PostDownloadedValue)
	InfoField := container.NewVBox(InfoTitle, InfoForm)
	//Work
	WorkLabel := widget.NewLabel("Work")
	WorkTitle := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), WorkLabel, layout.NewSpacer())
	WorkersLabel := widget.NewLabel("Workers")
	WorkersChoice := widget.NewSelect([]string{"1", "2", "3"}, nil)
	LimitLabel := widget.NewLabel("Limits")
	LimitValue := widget.NewEntry()
	LimitValue.SetPlaceHolder("Enter Number or leave it empty")
	StartPage := widget.NewLabel("Start Page")
	StartPageValue := widget.NewEntry()
	JobForm := container.New(layout.NewFormLayout(), WorkersLabel, WorkersChoice, LimitLabel, LimitValue, StartPage, StartPageValue)
	StartBtn := widget.NewButtonWithIcon("start", theme.DownloadIcon(), func() {
		log.Println(blog)
		StartScraping(blog)
	})
	BtnLayout := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), StartBtn, layout.NewSpacer())
	//ProgressBar
	progressBar := widget.NewProgressBar()
	Jobfield := container.New(layout.NewVBoxLayout(), WorkTitle, JobForm, BtnLayout, progressBar)

	mainContent := container.NewVSplit(InfoField, Jobfield)
	w.SetContent(mainContent)
	return w

}
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
		value2.Text = fmt.Sprint(len(post.Images))
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
