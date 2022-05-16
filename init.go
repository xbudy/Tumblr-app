package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func InitBlog(blog string) {
	log.Println("initializing ..")
	userpath := filepath.Join("blogs", blog)
	os.MkdirAll(userpath, os.ModePerm)

	log.Println("done")
}
