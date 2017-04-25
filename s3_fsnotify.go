package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fsnotify/fsnotify"
)

var contentTypes = map[string]string{
	"txt":  "text/plain",
	"json": "application/json",
	"xml":  "application/xml",
	"pdf":  "application/pdf",
	"html": "text/html",
	"htm":  "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"bmp":  "image/bmp",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"tiff": "image/tiff",
	"gif":  "image/gif",
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	bucket := os.Getenv("bucket")
	if bucket == "" {
		log.Fatalln("bucket env var must be set")
	}
	syncdir := os.Getenv("syncdir")
	if syncdir == "" {
		syncdir = "."
	}

	if _, err := os.Stat(syncdir); os.IsNotExist(err) {
		os.MkdirAll(syncdir, os.ModePerm)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create && !excludeFileExt(event.Name) {
					log.Println("modified file:", event.Name)
					file, err := os.Open(event.Name)
					if err != nil {
						log.Fatalln("Unable to open file", err)
					}
					defer file.Close()
					filename := path.Base(event.Name)
					sess := session.Must(session.NewSessionWithOptions(session.Options{
						SharedConfigState: session.SharedConfigEnable,
					}))
					uploader := s3manager.NewUploader(sess)
					_, err = uploader.Upload(&s3manager.UploadInput{
						Bucket:               aws.String(bucket),
						Key:                  aws.String(filename),
						Body:                 file,
						ServerSideEncryption: aws.String("AES256"),
						ContentType:          aws.String(getContentType(file)),
					})
					if err != nil {
						// Print the error and exit.
						log.Fatalln("Unable to upload", event.Name, bucket, err)
					}
					fmt.Println("Successfully uploaded", event.Name, "to", bucket)
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(syncdir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func getContentType(file *os.File) string {
	result := "binary/octet-stream"
	name := file.Name()
	pos := strings.LastIndex(name, ".")
	if pos != -1 {
		v, ok := contentTypes[name[pos+1:]]
		if ok {
			result = v
		}
	}
	return result
}

func excludeFileExt(path string) bool {
	var extension = filepath.Ext(path)
	switch extension {
	case ".swp":
		return true
	case ".tmp":
		return true
	default:
		return false
	}
	return false
}
