package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/fsnotify/fsnotify"
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
	filename := os.Getenv("filename")
	if filename == "" {
		log.Fatalln("filename env var must be set")
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					file, err := os.Open(filename)
					if err != nil {
						log.Fatalln("Unable to open file", err)
					}
					defer file.Close()

					uploader := s3manager.NewUploader(sess)

					_, err = uploader.Upload(&s3manager.UploadInput{
						Bucket:               aws.String(bucket),
						Key:                  aws.String(filename),
						Body:                 file,
						ServerSideEncryption: aws.String("AES256"),
					})
					if err != nil {
						// Print the error and exit.
						log.Fatalln("Unable to upload", filename, bucket, err)
					}
					fmt.Printf("Successfully uploaded", filename, bucket)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
