package main

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fsnotify/fsnotify"
)

const (
	FileExtensionSwp = ".swp"
	FileExtensionTmp = ".tmp"
)

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
	go watchEvents(watcher, bucket)
	err = watcher.Add(syncdir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func watchEvents(watcher *fsnotify.Watcher, bucket string) {
	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			if (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) && !excludeFileExt(event.Name) {
				uploadFileToS3(event.Name, bucket)
			}

		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func uploadFileToS3(filePath string, bucket string) {
	log.Println("modified file:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Unable to open file", err)
	}
	defer file.Close()

	filename := path.Base(filePath)
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
	file.Close()

	if err != nil {
		log.Fatalln("Unable to upload", filePath, bucket, err)
	}
	fmt.Println("Successfully uploaded", filePath, "to", bucket)
}

func getContentType(file *os.File) string {
	ext := filepath.Ext(file.Name())
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func excludeFileExt(path string) bool {
	extension := filepath.Ext(path)
	return extension == FileExtensionSwp || extension == FileExtensionTmp
}
