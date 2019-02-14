package models

import (
	     "github.com/jinzhu/gorm"
      "net/http"
      "io/ioutil"
      "lightupon-api/services/aws"
      // "lightupon-api/services/imageMagick"
      "fmt"
)

type Pin struct {
	gorm.Model
  Url string
  PostID uint
}

func NewPin(url string, id uint) (pin Pin, err error) {
	// binary, err := DownloadImage(url); if err != nil {
 //    return
 //  }

	// pinBinary := imageMagick.CropPin(binary)
  pinBinary := []byte{}

  asset := aws.Asset{
    Type: "images", 
    Name: assetName(id), 
    Extension: ".png",
    Binary: pinBinary,
  }

  getUrl, err := aws.UploadAsset(asset); if err != nil {
  	return
  }

  pin = Pin{
  	PostID: id,
  	Url: getUrl,
  }

  err = DB.Create(&pin).Error

  return
}

func assetName(id uint) string {
  return "/posts/" + fmt.Sprintf("%v", id) + "/pin"
}

func DownloadImage(url string) (imageBinary []byte, err error) {
  resp, err := http.Get(url); if err != nil {
    return
  }

  defer resp.Body.Close()

  imageBinary, err = ioutil.ReadAll(resp.Body); if err != nil {
    return
  }

  return
}
