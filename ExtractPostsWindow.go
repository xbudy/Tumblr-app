package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateExtractPostsWindow(mdb MongoDb, a fyne.App, blog string) fyne.Window {
	//Context of running
	ctx := context.WithValue(context.Background(), "run", true)
	stopped := true
	results := make(chan bool)
	w := a.NewWindow("Extract " + blog + " Posts")
	w.SetCloseIntercept(func() {
		log.Println("closing")
		if !stopped {
			results <- false
			close(results)
		}

		w.Close()
	})
	w.Resize(fyne.NewSize(780, 400))
	infoLabel := widget.NewLabel("Info")
	InfoTitle := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), infoLabel, layout.NewSpacer())
	LastPageLabel := widget.NewLabel("Last Page")
	LastPageValue := widget.NewLabel(mdb.GetLastPage(blog))
	LastPageValue.Wrapping = fyne.TextTruncate
	TotalPostsOnTLabel := widget.NewLabel("Posts On Tumblr")
	TotalPostsOnTValue := widget.NewLabel(fmt.Sprint(mdb.GetBlogInfo(blog).TotalPosts))
	PostDownloadedLabel := widget.NewLabel("Downloaded Posts")
	postsCount, _ := mdb.GetPostsCount(blog)
	PostDownloadedValue := widget.NewLabel(fmt.Sprint(postsCount))
	InfoForm := container.New(layout.NewFormLayout(), LastPageLabel, LastPageValue, PostDownloadedLabel, PostDownloadedValue, TotalPostsOnTLabel, TotalPostsOnTValue)
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
	//ProgressBar
	progressBar := widget.NewProgressBar()
	StartBtn := widget.NewButtonWithIcon("start", theme.DownloadIcon(), func() {
		log.Println(blog)
		if stopped {
			go StartScraping(results, mdb, blog, true, progressBar)
			stopped = false

		}

	})
	DwnBtn := widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
		InitBlog(blog)
		Pooling(blog, mdb.GetPosts(blog))

	})
	StopBtn := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		if !stopped {
			results <- false
			close(results)
			log.Println("Stopping")
			stopped = true
		}
		log.Println(ctx)

	})
	BtnLayout := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), StartBtn, DwnBtn, StopBtn, layout.NewSpacer())

	Jobfield := container.New(layout.NewVBoxLayout(), WorkTitle, JobForm, BtnLayout, progressBar)

	mainContent := container.NewVSplit(InfoField, Jobfield)
	w.SetContent(mainContent)
	return w

}
func newpBool(b bool) *bool {
	return &b
}
