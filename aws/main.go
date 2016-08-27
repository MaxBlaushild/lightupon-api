package aws

import (
      "github.com/aws/aws-sdk-go/aws"
      "github.com/aws/aws-sdk-go/aws/session"
      "github.com/aws/aws-sdk-go/service/s3"
      "time"
      )

const AUDIO string = "audio"
const VIDEO string = "video"
const IMAGE string = "image"

var UPLOAD_TYPES [3]string = [3]string{ AUDIO, VIDEO, IMAGE }

func PutAsset(assetType string, assetName string)(urlStr string, err error) {
	LightuponS3 := startS3Session()
	key := formKey(assetType, assetName)

	req, _ := LightuponS3.PutObjectRequest(&s3.PutObjectInput{
	  Bucket: aws.String("lightupon"),
	  Key:    aws.String(key),
	})

	urlStr, err = req.Presign(15 * time.Minute)
	return
}

func GetAsset(assetType string, assetName string) (urlStr string, err error) {
	LightuponS3 := startS3Session()
	key := formKey(assetType, assetName)

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

func formKey(assetType string, assetName string) string {
	return assetType + "/" + assetName
}
