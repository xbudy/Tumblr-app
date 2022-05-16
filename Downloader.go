package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func DownloadFromUrl(blog string, Timestamp int, Id string, url string, i int, overwrite bool) error {
	path, _ := os.Getwd()
	newpath := filepath.Join(path)
	filename := blog + "-" + fmt.Sprint(i) + "_" + fmt.Sprint(Timestamp) + "_" + Id
	filename = fmt.Sprint(filename) + filepath.Ext(url)
	fileName := newpath + "/blogs/" + blog + "/" + filename
	var err error

	if _, err := os.Stat(fileName); err == nil && !overwrite {
		log.Println("exist")
		return nil

	} else if errors.Is(err, os.ErrNotExist) || err == nil && overwrite {
		// path/to/whatever does *not* exist

		resp, e := http.Get(url)
		if e != nil {
			log.Error("Download:request error", e)
		} else if e == nil && strings.Contains(resp.Header.Get("Content-Type"), "image") {
			defer resp.Body.Close()
			//Create new file for images
			file, err := os.Create(fileName)
			if err == nil {
				file.ReadFrom(resp.Body)
				log.Info("Downloaded !:", filename)
			} else {
				log.Error("Error while creating file")
				return err
			}
			defer file.Close()
		}
	}
	return err
}
