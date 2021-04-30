package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3session *s3.S3

const (
	REGION      = "eu-west-1"
	BUCKET_NAME = "lambdatestyoutube"
)

type InputEvent struct {
	Link string `json:"link"`
	Key  string `json:"key"`
}

func init() {
	// create a new S3 Session
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func main() {
	// start the lambda handler
	lambda.Start(Handler)
}

func Handler(event InputEvent) (int, error) {
	// get the image
	image := GetImage(event.Link)

	// upload it to s3 using PutObject
	_, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(image),
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(event.Key),
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
}

func GetImage(url string) (bytes []byte) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	return bytes
}
