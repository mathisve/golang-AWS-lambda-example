package main

import (
  "fmt"
  "github.com/aws/aws-lambda-go/lambda"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"

  "net/http"
  "io/ioutil"
  "bytes"
)

var s3session *s3.S3

const (
  REGION = "eu-west-1"
  BUCKET_NAME = "lambdatestyoutube"
)

func init() {
  s3session = s3.New(session.Must(session.NewSession(&aws.Config{
      Region: aws.String(REGION),
    })))
}

func main() {
  lambda.Start(Handler)
}

func Handler(event InputEvent) (string, error) {
  image := GetImage(event.Link)
  _, err := s3session.PutObject(&s3.PutObjectInput{
    Body: bytes.NewReader(image),
    Bucket: aws.String(BUCKET_NAME),
    Key: aws.String(event.Key),
  })

  if err != nil {
    return "Something went wrong!", err
  }

  return "Everything worked YAY!", err
}

type InputEvent struct {
  Link string `json:"link"`
  Key  string `json:"key"`
}

func GetImage(url string) (bytes []byte) {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
  }
  defer resp.Body.Close()

  bytes, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
  }

  return bytes
}
