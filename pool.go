package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Pooling(blog string, Jobs []PostMeta) {
	jobs := make(chan PostMeta, len(Jobs))
	JobsResults := make(chan string, len(Jobs))

	for w := 1; w <= 3; w++ {
		go worker(blog, w, jobs, JobsResults)
	}
	for j := 0; j < len(Jobs); j++ {
		jobs <- Jobs[j]
	}
	close(jobs)
	for a := 1; a <= len(Jobs); a++ {
		<-JobsResults
	}
}

func worker(blog string, id int, jobs <-chan PostMeta, JobsResults chan<- string) {
	for j := range jobs {
		log.Println("worker", id, "started  job", j.Timestamp)
		for i, img := range j.Medias {
			log.Println(img)
			err := DownloadFromUrl(blog, j.Timestamp, j.Id, img.Url, i, false)
			if err != nil {
				log.Error(err)
			}
		}
		log.Println("worker", id, "finished job", j.Timestamp)
		JobsResults <- fmt.Sprint(j.Timestamp)
	}
}
