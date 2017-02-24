package aws

import (
      "github.com/aws/aws-sdk-go/aws"
      "github.com/aws/aws-sdk-go/aws/session"
      "github.com/aws/aws-sdk-go/service/s3"
      "time"
      "strings"
      "net/http"
      "bytes"
      "fmt"
      )

type Asset struct {
	Type string
	Name string
	Extension string
	Binary []byte
}

func PutAsset(asset Asset)(urlStr string, err error) {
	LightuponS3 := startS3Session()
	key := formKey(asset)

	req, _ := LightuponS3.PutObjectRequest(&s3.PutObjectInput{
	  Bucket: aws.String("lightupon"),
	  Key:    aws.String(key),
	  ACL:    aws.String("public-read"),
	})

	urlStr, err = req.Presign(15 * time.Minute)
	return
}

func UploadAsset(asset Asset) (getUrl string, err error) {
  putUrl, err := PutAsset(asset)
  client := &http.Client{}
  request, err := http.NewRequest("PUT", putUrl, bytes.NewReader(asset.Binary))
  request.Header.Set("Content-Type", "image/png")
  request.Header.Set("x-amz-acl", "public-read")
  response, err := client.Do(request)
  defer response.Body.Close()
 	fmt.Println(response.Status)
  if err == nil && response.StatusCode == 200 {
    getUrl = splitPutUrl(putUrl)
  }
  return
}

func splitPutUrl(putUrl string) (getUrl string){
	urlSegments := strings.Split(putUrl, "?")
	getUrl = urlSegments[0]
	return
}

func GetAsset(asset Asset) (urlStr string, err error) {
	LightuponS3 := startS3Session()
	key := formKey(asset)

	req, _ := LightuponS3.GetObjectRequest(&s3.GetObjectInput{
	  Bucket: aws.String("lightupon"),
	  Key:    aws.String(key),
	})

	urlStr, err = req.Presign(15 * time.Minute)
	return
}

func startS3Session() *s3.S3 {
	s3 := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	return s3
}

func formKey(asset Asset) string {
	return asset.Type + "/" + asset.Name + asset.Extension
}
